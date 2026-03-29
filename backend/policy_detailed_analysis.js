const puppeteer = require('puppeteer');
const fs = require('fs');

(async () => {
  console.log('=== 启动浏览器 ===');
  const browser = await puppeteer.launch({
    headless: true,
    args: [
      '--no-sandbox',
      '--disable-setuid-sandbox',
      '--disable-dev-shm-usage',
      '--window-size=1920,1080',
      '--ignore-certificate-errors',
      '--ignore-ssl-errors',
      '--disable-web-security'
    ]
  });

  const page = await browser.newPage();
  await page.setViewport({ width: 1920, height: 1080 });
  
  console.log('=== 访问登录页面 ===');
  await page.goto('https://192.168.31.39/auth/login', { 
    waitUntil: 'networkidle2',
    timeout: 60000 
  });
  
  console.log('正在登录...');
  await page.waitForSelector('input[placeholder*="用户名"]', { timeout: 10000 });
  await page.click('input[placeholder*="用户名"]');
  await page.type('input[placeholder*="用户名"]', 'admin', { delay: 50 });
  
  await page.click('input[placeholder*="密码"]');
  await page.type('input[placeholder*="密码"]', 'admin@123', { delay: 50 });
  
  await page.waitForSelector('button.ant-btn-primary', { timeout: 5000 });
  await page.click('button.ant-btn-primary');
  
  console.log('等待登录完成...');
  await new Promise(r => setTimeout(r, 5000));
  await page.waitForNavigation({ waitUntil: 'networkidle2', timeout: 15000 }).catch(() => {});
  
  console.log('=== 访问政策页面 ===');
  await page.goto('https://192.168.31.39/policy', { 
    waitUntil: 'networkidle2',
    timeout: 60000 
  });
  
  console.log('等待页面加载...');
  await new Promise(r => setTimeout(r, 5000));
  
  console.log('\n=== 详细页面分析 ===\n');
  
  const detailedAnalysis = await page.evaluate(() => {
    const result = {
      title: document.title,
      url: window.location.href,
      
      // 页面整体结构
      layout: {
        hasSidebar: !!document.querySelector('.sidebar, .side-menu, .ant-menu, [class*="sidebar"], [class*="side-menu"]'),
        hasHeader: !!document.querySelector('header, .header, .top-nav, [class*="header"], [class*="top-nav"]'),
        hasTabs: !!document.querySelector('.tabs, .tab, [role="tablist"], [class*="tab"]'),
        sidebarText: '',
        headerText: ''
      },
      
      // 权限列表信息
      permissionList: {
        hasTable: false,
        tableHeaders: [],
        totalRows: 0,
        permissions: []
      },
      
      // 权限分类
      categories: {
        hasCategories: false,
        categoryTabs: [],
        currentCategory: ''
      },
      
      // 搜索和过滤
      searchFilter: {
        hasSearch: false,
        searchPlaceholder: '',
        hasFilters: false,
        filterOptions: []
      },
      
      // 操作按钮
      actions: {
        hasCreate: false,
        hasDelete: false,
        hasEdit: false,
        hasEnable: false,
        hasDisable: false,
        hasShare: false,
        buttons: []
      },
      
      // 分页信息
      pagination: {
        hasPagination: false,
        currentPage: '',
        totalPages: '',
        totalRecords: '',
        pageSize: ''
      }
    };
    
    // 分析布局
    const sidebar = document.querySelector('.sidebar, .side-menu, .ant-menu-vertical, [class*="sidebar"], [class*="side-menu"]');
    if (sidebar) {
      result.layout.sidebarText = sidebar.innerText.substring(0, 500);
    }
    
    const header = document.querySelector('header, .header, .top-nav, [class*="header"]');
    if (header) {
      result.layout.headerText = header.innerText.substring(0, 200);
    }
    
    // 分析表格
    const tables = document.querySelectorAll('table, .vxe-table, [class*="table"]');
    if (tables.length > 0) {
      result.permissionList.hasTable = true;
      
      // 获取表头
      const mainTable = tables[tables.length - 1]; // 最后一个表格通常是数据表格
      const headers = mainTable.querySelectorAll('thead th, .vxe-table--header th, [class*="header"] th');
      headers.forEach(th => {
        const text = th.innerText.trim();
        if (text) result.permissionList.tableHeaders.push(text);
      });
      
      // 获取所有权限行
      const rows = mainTable.querySelectorAll('tbody tr, .vxe-table--body tr, [class*="body"] tr');
      result.permissionList.totalRows = rows.length;
      
      rows.forEach((row, index) => {
        if (index < 50) { // 只获取前 50 行
          const cells = row.querySelectorAll('td');
          const rowData = {};
          result.permissionList.tableHeaders.forEach((header, i) => {
            rowData[header] = cells[i]?.innerText.trim() || '';
          });
          result.permissionList.permissions.push(rowData);
        }
      });
    }
    
    // 分析分类标签
    const tabs = document.querySelectorAll('.tabs .ant-tabs-tab, [class*="tab"]:not([class*="table"])');
    if (tabs.length > 0) {
      result.categories.hasCategories = true;
      tabs.forEach(tab => {
        const text = tab.innerText.trim();
        if (text && !result.categories.categoryTabs.includes(text)) {
          result.categories.categoryTabs.push(text);
        }
      });
    }
    
    // 查找分类选择器
    const categorySelectors = document.querySelectorAll('.category-select, .filter-select, .ant-select');
    categorySelectors.forEach(sel => {
      const text = sel.innerText.trim();
      if (text && (text.includes('权限') || text.includes('系统') || text.includes('自定义'))) {
        result.categories.hasCategories = true;
        if (!result.categories.categoryTabs.includes(text)) {
          result.categories.categoryTabs.push(text);
        }
      }
    });
    
    // 分析搜索框
    const searchInput = document.querySelector('input[placeholder*="搜索"], input[placeholder*="Search"], .search-input');
    if (searchInput) {
      result.searchFilter.hasSearch = true;
      result.searchFilter.searchPlaceholder = searchInput.placeholder;
    }
    
    // 分析过滤选项
    const selects = document.querySelectorAll('select, .ant-select');
    selects.forEach(sel => {
      const text = sel.innerText.trim();
      if (text && text.length < 100) {
        result.searchFilter.hasFilters = true;
        result.searchFilter.filterOptions.push(text);
      }
    });
    
    // 分析操作按钮
    const allButtons = document.querySelectorAll('button, .ant-btn, [role="button"]');
    allButtons.forEach(btn => {
      const text = btn.innerText.trim();
      if (text) {
        result.actions.buttons.push(text);
        if (text.includes('新建') || text.includes('创建') || text.includes('Add') || text.includes('Create')) {
          result.actions.hasCreate = true;
        }
        if (text.includes('删除') || text.includes('Delete')) {
          result.actions.hasDelete = true;
        }
        if (text.includes('修改') || text.includes('编辑') || text.includes('Edit')) {
          result.actions.hasEdit = true;
        }
        if (text.includes('启用') || text.includes('Enable')) {
          result.actions.hasEnable = true;
        }
        if (text.includes('禁用') || text.includes('Disable')) {
          result.actions.hasDisable = true;
        }
        if (text.includes('共享') || text.includes('Share')) {
          result.actions.hasShare = true;
        }
      }
    });
    
    // 分析分页
    const pagination = document.querySelector('.pagination, .ant-pagination, [class*="pagination"]');
    if (pagination) {
      result.pagination.hasPagination = true;
      const paginationText = pagination.innerText;
      
      // 提取总记录数
      const totalMatch = paginationText.match(/共\s*(\d+)\s*条/);
      if (totalMatch) result.pagination.totalRecords = totalMatch[1];
      
      // 提取当前页
      const pageMatch = paginationText.match(/第\s*(\d+)\s*页/);
      if (pageMatch) result.pagination.currentPage = pageMatch[1];
    }
    
    return result;
  });
  
  // 输出详细分析结果
  console.log('=== 页面基本信息 ===');
  console.log(`标题：${detailedAnalysis.title}`);
  console.log(`URL: ${detailedAnalysis.url}`);
  
  console.log('\n=== 页面整体结构和布局 ===');
  console.log(`侧边栏：${detailedAnalysis.layout.hasSidebar ? '有' : '无'}`);
  console.log(`顶部导航：${detailedAnalysis.layout.hasHeader ? '有' : '无'}`);
  console.log(`标签页：${detailedAnalysis.layout.hasTabs ? '有' : '无'}`);
  
  if (detailedAnalysis.layout.sidebarText) {
    console.log('\n侧边栏菜单项:');
    detailedAnalysis.layout.sidebarText.split('\n').filter(t => t.trim()).forEach(t => {
      console.log(`  - ${t.trim()}`);
    });
  }
  
  console.log('\n=== 权限列表展示方式 ===');
  console.log(`使用表格展示：${detailedAnalysis.permissionList.hasTable ? '是' : '否'}`);
  console.log(`表格列头：${detailedAnalysis.permissionList.tableHeaders.join(' | ')}`);
  console.log(`默认展示数量：${detailedAnalysis.permissionList.totalRows} 条`);
  
  console.log('\n=== 权限分类 ===');
  console.log(`有分类功能：${detailedAnalysis.categories.hasCategories ? '是' : '否'}`);
  if (detailedAnalysis.categories.categoryTabs.length > 0) {
    console.log('分类标签:');
    detailedAnalysis.categories.categoryTabs.forEach(cat => {
      console.log(`  - ${cat}`);
    });
  }
  
  console.log('\n=== 权限具体字段和格式 ===');
  if (detailedAnalysis.permissionList.permissions.length > 0) {
    console.log('字段说明:');
    detailedAnalysis.permissionList.tableHeaders.forEach(header => {
      console.log(`  - ${header}`);
    });
    
    console.log('\n前 10 条权限示例:');
    detailedAnalysis.permissionList.permissions.slice(0, 10).forEach((perm, i) => {
      console.log(`\n${i + 1}. ${perm['名称'] || perm['name'] || '未知'}`);
      Object.entries(perm).forEach(([key, value]) => {
        if (key !== '名称' && key !== 'name' && value) {
          console.log(`   ${key}: ${value}`);
        }
      });
    });
  }
  
  console.log('\n=== 搜索和过滤功能 ===');
  console.log(`搜索功能：${detailedAnalysis.searchFilter.hasSearch ? '有' : '无'}`);
  if (detailedAnalysis.searchFilter.searchPlaceholder) {
    console.log(`搜索提示：${detailedAnalysis.searchFilter.searchPlaceholder}`);
  }
  console.log(`过滤功能：${detailedAnalysis.searchFilter.hasFilters ? '有' : '无'}`);
  if (detailedAnalysis.searchFilter.filterOptions.length > 0) {
    console.log('过滤选项:');
    detailedAnalysis.searchFilter.filterOptions.slice(0, 10).forEach(opt => {
      console.log(`  - ${opt}`);
    });
  }
  
  console.log('\n=== 操作按钮 ===');
  console.log(`新建：${detailedAnalysis.actions.hasCreate ? '有' : '无'}`);
  console.log(`删除：${detailedAnalysis.actions.hasDelete ? '有' : '无'}`);
  console.log(`修改：${detailedAnalysis.actions.hasEdit ? '有' : '无'}`);
  console.log(`启用：${detailedAnalysis.actions.hasEnable ? '有' : '无'}`);
  console.log(`禁用：${detailedAnalysis.actions.hasDisable ? '有' : '无'}`);
  console.log(`共享：${detailedAnalysis.actions.hasShare ? '有' : '无'}`);
  
  console.log('\n=== 分页信息 ===');
  console.log(`有分页：${detailedAnalysis.pagination.hasPagination ? '是' : '否'}`);
  console.log(`总记录数：${detailedAnalysis.pagination.totalRecords || '未知'}`);
  
  // 保存完整数据
  fs.writeFileSync(
    '/Users/aurorapeng/Desktop/work/git/new/openCMP/backend/policy_detailed_analysis.json',
    JSON.stringify(detailedAnalysis, null, 2)
  );
  console.log('\n完整分析数据已保存到：policy_detailed_analysis.json');
  
  // 截图
  await page.screenshot({ 
    path: '/Users/aurorapeng/Desktop/work/git/new/openCMP/backend/policy_page_full.png',
    fullPage: true 
  });
  console.log('完整截图已保存到：policy_page_full.png');
  
  await browser.close();
  console.log('\n=== 完成 ===');
})().catch(err => {
  console.error('错误:', err);
  process.exit(1);
});
