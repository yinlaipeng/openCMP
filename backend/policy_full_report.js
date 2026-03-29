const puppeteer = require('puppeteer');
const fs = require('fs');

(async () => {
  console.log('=== 启动浏览器并访问权限页面 ===');
  const browser = await puppeteer.launch({
    headless: true,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--window-size=1920,1080',
      '--ignore-certificate-errors',
      '--ignore-ssl-errors'
    ]
  });

  const page = await browser.newPage();
  await page.setViewport({ width: 1920, height: 1080 });
  
  // 登录
  console.log('正在登录...');
  await page.goto('https://192.168.31.39/auth/login', { waitUntil: 'networkidle2', timeout: 60000 });
  await page.waitForSelector('input[placeholder*="用户名"]', { timeout: 10000 });
  await page.click('input[placeholder*="用户名"]');
  await page.type('input[placeholder*="用户名"]', 'admin', { delay: 50 });
  await page.click('input[placeholder*="密码"]');
  await page.type('input[placeholder*="密码"]', 'admin@123', { delay: 50 });
  await page.click('button.ant-btn-primary');
  await new Promise(r => setTimeout(r, 5000));
  await page.waitForNavigation({ waitUntil: 'networkidle2', timeout: 15000 }).catch(() => {});
  
  // 访问政策页面
  console.log('访问权限页面...');
  await page.goto('https://192.168.31.39/policy', { waitUntil: 'networkidle2', timeout: 60000 });
  await new Promise(r => setTimeout(r, 5000));
  
  console.log('\n=== 完整页面分析报告 ===\n');
  
  const fullAnalysis = await page.evaluate(() => {
    const result = {
      pageTitle: document.title,
      pageUrl: window.location.href,
      
      // 1. 页面整体结构和布局
      layout: {
        description: '',
        mainComponents: []
      },
      
      // 2. 权限列表展示方式
      permissionTable: {
        type: '',
        headers: [],
        rowCount: 0,
        data: []
      },
      
      // 3. 权限分类
      categories: {
        hasCategories: false,
        tabs: [],
        currentTab: ''
      },
      
      // 4. 权限具体字段和格式
      permissionFields: [],
      permissionExamples: [],
      
      // 5. 搜索、过滤功能
      searchFilter: {
        hasSearch: false,
        searchInfo: '',
        hasFilters: false,
        filterInfo: ''
      },
      
      // 6. 默认展示数量
      displayInfo: {
        defaultCount: 0,
        totalCount: 0,
        hasPagination: false,
        paginationInfo: ''
      },
      
      // 操作按钮
      actions: []
    };
    
    // 获取页面所有文本
    const allText = document.body.innerText;
    
    // 分析布局
    const app = document.getElementById('app');
    if (app) {
      const header = app.querySelector('.header, header, .top-nav, .navbar')?.innerText?.substring(0, 200);
      const main = app.querySelector('.main, main, .content, .page-content')?.innerText?.substring(0, 500);
      
      result.layout.mainComponents.push('顶部导航栏');
      result.layout.mainComponents.push('主内容区');
      result.layout.description = '页面采用典型的后台管理系统布局，包含顶部导航和主内容区域';
    }
    
    // 分析表格 - 尝试多种方式获取表格数据
    const tables = document.querySelectorAll('table, .vxe-table, .ant-table, [class*="table"]');
    
    // 尝试从 vxe-table 获取数据
    const vxeTable = document.querySelector('.vxe-table');
    if (vxeTable) {
      result.permissionTable.type = 'vxe-table';
      
      // 获取表头
      const headerCells = vxeTable.querySelectorAll('.vxe-table--header th, .vxe-table--header .vxe-cell');
      headerCells.forEach(cell => {
        const text = cell.innerText.trim();
        if (text) result.permissionTable.headers.push(text);
      });
      
      // 获取数据行
      const bodyRows = vxeTable.querySelectorAll('.vxe-table--body tr');
      result.permissionTable.rowCount = bodyRows.length;
      
      bodyRows.forEach((row, idx) => {
        if (idx < 30) {
          const cells = row.querySelectorAll('td, .vxe-cell');
          const rowData = {};
          cells.forEach((cell, cellIdx) => {
            const header = result.permissionTable.headers[cellIdx] || `列${cellIdx + 1}`;
            rowData[header] = cell.innerText.trim();
          });
          result.permissionTable.data.push(rowData);
        }
      });
    }
    
    // 如果 vxe-table 没有数据，尝试 ant-table
    if (result.permissionTable.data.length === 0) {
      const antTable = document.querySelector('.ant-table');
      if (antTable) {
        result.permissionTable.type = 'ant-table';
        
        const headerCells = antTable.querySelectorAll('.ant-table-thead th');
        headerCells.forEach(cell => {
          const text = cell.innerText.trim();
          if (text) result.permissionTable.headers.push(text);
        });
        
        const bodyRows = antTable.querySelectorAll('.ant-table-tbody tr');
        result.permissionTable.rowCount = bodyRows.length;
        
        bodyRows.forEach((row, idx) => {
          if (idx < 30) {
            const cells = row.querySelectorAll('td');
            const rowData = {};
            cells.forEach((cell, cellIdx) => {
              const header = result.permissionTable.headers[cellIdx] || `列${cellIdx + 1}`;
              rowData[header] = cell.innerText.trim();
            });
            result.permissionTable.data.push(rowData);
          }
        });
      }
    }
    
    // 分析分类标签
    const tabs = document.querySelectorAll('.ant-tabs-tab, .tabs .tab, [role="tab"]');
    if (tabs.length > 0) {
      result.categories.hasCategories = true;
      tabs.forEach(tab => {
        const text = tab.innerText.trim();
        if (text && !result.categories.tabs.includes(text)) {
          result.categories.tabs.push(text);
        }
      });
    }
    
    // 查找分类文本（可能是静态文本）
    const categoryText = allText.match(/(全部 | 自定义权限 | 系统权限)/g);
    if (categoryText) {
      result.categories.hasCategories = true;
      const uniqueCats = [...new Set(categoryText)];
      uniqueCats.forEach(cat => {
        if (!result.categories.tabs.includes(cat)) {
          result.categories.tabs.push(cat);
        }
      });
    }
    
    // 分析搜索框
    const searchInput = document.querySelector('input[placeholder*="搜索"], input[placeholder*="Search"]');
    if (searchInput) {
      result.searchFilter.hasSearch = true;
      result.searchFilter.searchInfo = searchInput.placeholder;
    }
    
    // 查找搜索提示文本
    const searchHint = allText.match(/默认为名称搜索.*?/);
    if (searchHint) {
      result.searchFilter.hasSearch = true;
      result.searchFilter.searchInfo = searchHint[0];
    }
    
    // 分析字段
    result.permissionFields = result.permissionTable.headers;
    
    // 获取权限示例
    result.permissionExamples = result.permissionTable.data.slice(0, 10);
    
    // 分析操作按钮
    const buttons = document.querySelectorAll('button, .ant-btn, [role="button"]');
    buttons.forEach(btn => {
      const text = btn.innerText.trim();
      if (text && !result.actions.includes(text)) {
        result.actions.push(text);
      }
    });
    
    // 分析分页
    const pagination = document.querySelector('.ant-pagination, .pagination');
    if (pagination) {
      result.displayInfo.hasPagination = true;
      result.displayInfo.paginationInfo = pagination.innerText;
      
      const totalMatch = pagination.innerText.match(/共\s*(\d+)\s*条/);
      if (totalMatch) {
        result.displayInfo.totalCount = parseInt(totalMatch[1]);
      }
    }
    
    // 如果没有分页信息，从文本中查找
    const totalMatch = allText.match(/共\s*(\d+)\s*条/);
    if (totalMatch) {
      result.displayInfo.totalCount = parseInt(totalMatch[1]);
      result.displayInfo.hasPagination = true;
    }
    
    // 默认展示数量
    result.displayInfo.defaultCount = result.permissionTable.rowCount;
    
    return result;
  });
  
  // 输出详细报告
  console.log('1. 页面整体结构和布局');
  console.log('   ====================');
  console.log(`   页面标题：${fullAnalysis.pageTitle}`);
  console.log(`   页面 URL: ${fullAnalysis.pageUrl}`);
  console.log(`   布局描述：${fullAnalysis.layout.description}`);
  console.log(`   主要组件：${fullAnalysis.layout.mainComponents.join(', ')}`);
  console.log('');
  
  console.log('2. 权限列表展示方式');
  console.log('   ====================');
  console.log(`   表格类型：${fullAnalysis.permissionTable.type || '未检测到'}`);
  console.log(`   表格列头：${fullAnalysis.permissionTable.headers.join(' | ')}`);
  console.log(`   当前显示行数：${fullAnalysis.permissionTable.rowCount}`);
  console.log('');
  
  console.log('3. 权限分类');
  console.log('   ====================');
  console.log(`   是否有分类：${fullAnalysis.categories.hasCategories ? '是' : '否'}`);
  if (fullAnalysis.categories.tabs.length > 0) {
    console.log('   分类标签:');
    fullAnalysis.categories.tabs.forEach(tab => {
      console.log(`     - ${tab}`);
    });
  }
  console.log('');
  
  console.log('4. 权限具体字段和格式');
  console.log('   ====================');
  console.log('   字段列表:');
  fullAnalysis.permissionFields.forEach((field, idx) => {
    console.log(`     ${idx + 1}. ${field}`);
  });
  console.log('');
  
  console.log('5. 权限示例数据（前 10 条）');
  console.log('   ====================');
  fullAnalysis.permissionExamples.forEach((perm, idx) => {
    console.log(`\n   ${idx + 1}. ${perm['名称'] || perm['name'] || Object.values(perm)[0] || '未知'}`);
    Object.entries(perm).forEach(([key, value]) => {
      if (key !== '名称' && key !== 'name' && key !== Object.keys(perm)[0] && value) {
        console.log(`      ${key}: ${value}`);
      }
    });
  });
  console.log('');
  
  console.log('6. 搜索、过滤功能');
  console.log('   ====================');
  console.log(`   搜索功能：${fullAnalysis.searchFilter.hasSearch ? '有' : '无'}`);
  if (fullAnalysis.searchFilter.searchInfo) {
    console.log(`   搜索说明：${fullAnalysis.searchFilter.searchInfo}`);
  }
  console.log(`   过滤功能：${fullAnalysis.searchFilter.hasFilters ? '有' : '无'}`);
  console.log('');
  
  console.log('7. 默认展示数量');
  console.log('   ====================');
  console.log(`   当前显示：${fullAnalysis.displayInfo.defaultCount} 条`);
  console.log(`   总记录数：${fullAnalysis.displayInfo.totalCount || '未知'} 条`);
  console.log(`   分页功能：${fullAnalysis.displayInfo.hasPagination ? '有' : '无'}`);
  if (fullAnalysis.displayInfo.paginationInfo) {
    console.log(`   分页信息：${fullAnalysis.displayInfo.paginationInfo}`);
  }
  console.log('');
  
  console.log('8. 操作按钮');
  console.log('   ====================');
  fullAnalysis.actions.forEach(action => {
    console.log(`   - ${action}`);
  });
  console.log('');
  
  // 保存完整数据
  fs.writeFileSync(
    '/Users/aurorapeng/Desktop/work/git/new/openCMP/backend/policy_full_report.json',
    JSON.stringify(fullAnalysis, null, 2)
  );
  console.log('完整报告已保存到：policy_full_report.json');
  
  await browser.close();
  console.log('\n=== 分析完成 ===');
})().catch(err => {
  console.error('错误:', err);
  process.exit(1);
});
