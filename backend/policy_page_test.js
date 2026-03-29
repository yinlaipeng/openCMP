const puppeteer = require('puppeteer');
const fs = require('fs');

(async () => {
  console.log('=== 步骤 1: 启动浏览器 ===');
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

  console.log('=== 步骤 2: 查看当前已打开的所有页面 ===');
  const pages = await browser.pages();
  console.log(`当前已打开 ${pages.length} 个页面`);

  console.log('\n=== 步骤 3: 创建新页面并访问 https://192.168.31.39/policy ===');
  const page = await browser.newPage();
  await page.setViewport({ width: 1920, height: 1080 });
  
  // 设置请求头
  await page.setExtraHTTPHeaders({
    'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8'
  });
  
  try {
    console.log('正在导航到页面...');
    await page.goto('https://192.168.31.39/policy', { 
      waitUntil: 'networkidle2',
      timeout: 60000 
    });
    console.log('页面加载成功');
  } catch (error) {
    console.log(`页面加载遇到问题：${error.message}`);
  }

  console.log('\n=== 步骤 4: 等待 5 秒确保内容完全加载 ===');
  await new Promise(r => setTimeout(r, 5000));
  
  console.log(`当前 URL: ${page.url()}`);
  const title = await page.title();
  console.log(`页面标题：${title}`);

  // 检查是否需要登录
  const loginCheck = await page.evaluate(() => {
    const bodyText = document.body.innerText;
    const hasLogin = bodyText.includes('登录') || bodyText.includes('login') || 
                     document.querySelector('input[type="password"]') !== null;
    const hasUsername = document.querySelector('input[name="username"], input[name="user"], #username, input[placeholder*="用户名"], input[placeholder*="user"]') !== null;
    return { hasLogin, hasUsername, bodyText: bodyText.substring(0, 500) };
  });
  
  console.log('登录检查结果:', loginCheck.hasLogin ? '需要登录' : '已登录或直接访问');
  
  if (loginCheck.hasLogin && loginCheck.hasUsername) {
    console.log('检测到登录页面，正在登录...');
    try {
      // 等待登录表单出现
      await page.waitForSelector('input[placeholder*="用户名"], input[name="username"], input[name="user"], #username', { timeout: 10000 });
      
      // 输入用户名
      const usernameSelector = 'input[placeholder*="用户名"], input[name="username"], input[name="user"], #username';
      await page.click(usernameSelector);
      await page.type(usernameSelector, 'admin', { delay: 50 });
      console.log('已输入用户名');
      
      // 输入密码
      const passwordSelector = 'input[placeholder*="密码"], input[name="password"], #password, input[type="password"]';
      await page.click(passwordSelector);
      await page.type(passwordSelector, 'admin@123', { delay: 50 });
      console.log('已输入密码');
      
      // 点击登录按钮
      const loginButton = 'button[type="submit"], input[type="submit"], .login-btn, #login-btn, button.ant-btn-primary, button:has-text("登 录"), button:has-text("登录")';
      await page.waitForSelector('button.ant-btn-primary', { timeout: 5000 });
      await page.click('button.ant-btn-primary');
      console.log('已点击登录按钮，等待跳转...');
      
      // 等待页面跳转
      await new Promise(r => setTimeout(r, 8000));
      await page.waitForNavigation({ waitUntil: 'networkidle2', timeout: 15000 }).catch(() => {
        console.log('导航超时，继续...');
      });
      
      console.log('当前 URL:', page.url());
    } catch (e) {
      console.log('登录操作遇到问题:', e.message);
    }
  }

  console.log('\n=== 步骤 5: 获取页面快照 ===');
  
  // 获取页面完整内容
  const htmlContent = await page.content();
  fs.writeFileSync('/Users/aurorapeng/Desktop/work/git/new/openCMP/backend/policy_page.html', htmlContent);
  console.log('HTML 已保存到：policy_page.html');

  // 截图
  await page.screenshot({ 
    path: '/Users/aurorapeng/Desktop/work/git/new/openCMP/backend/policy_page.png',
    fullPage: true 
  });
  console.log('截图已保存到：policy_page.png');

  // 详细分析页面内容
  console.log('\n=== 页面详细分析 ===');
  const pageAnalysis = await page.evaluate(() => {
    const result = {
      title: document.title,
      url: window.location.href,
      bodyText: document.body.innerText.substring(0, 5000),
      structure: [],
      tables: [],
      forms: [],
      inputs: [],
      buttons: [],
      links: [],
      permissionKeywords: []
    };

    // 顶层结构
    document.querySelectorAll('body > *').forEach((el, i) => {
      if (i < 20) {
        result.structure.push({
          tag: el.tagName.toLowerCase(),
          id: el.id || '',
          class: el.className || '',
          text: (el.innerText || '').substring(0, 100).replace(/\n/g, ' ')
        });
      }
    });

    // 表格
    document.querySelectorAll('table').forEach((table, i) => {
      const headers = [];
      table.querySelectorAll('thead th, thead td, tr:first-child th, tr:first-child td').forEach(th => {
        headers.push(th.innerText.trim());
      });
      const rows = table.querySelectorAll('tbody tr, tr');
      result.tables.push({
        index: i,
        headers: headers,
        rowCount: rows.length,
        firstRow: headers.length > 0 ? 
          Array.from(rows[0]?.querySelectorAll('td') || []).map(td => td.innerText.trim()) : []
      });
    });

    // 输入框
    document.querySelectorAll('input').forEach(input => {
      result.inputs.push({
        type: input.type,
        name: input.name,
        placeholder: input.placeholder,
        value: input.value
      });
    });

    // 按钮
    document.querySelectorAll('button, input[type="submit"], input[type="button"]').forEach(btn => {
      result.buttons.push({
        text: btn.innerText || btn.value,
        type: btn.type,
        class: btn.className
      });
    });

    // 查找权限相关关键词
    const allText = document.body.innerText;
    const keywords = ['权限', 'Permission', 'Policy', '策略', '角色', 'Role', '菜单', 'Menu'];
    keywords.forEach(kw => {
      const regex = new RegExp(`.{0,30}${kw}.{0,50}`, 'g');
      const matches = allText.match(regex);
      if (matches) {
        result.permissionKeywords.push({ keyword: kw, matches: matches.slice(0, 5) });
      }
    });

    return result;
  });

  console.log('\n--- 页面标题 ---');
  console.log(pageAnalysis.title);
  
  console.log('\n--- 页面 URL ---');
  console.log(pageAnalysis.url);

  console.log('\n--- 页面结构 ---');
  pageAnalysis.structure.forEach((el, i) => {
    console.log(`${i + 1}. <${el.tag}> ${el.id ? 'id="' + el.id + '"' : ''} ${el.class ? 'class="' + el.class + '"' : ''}`);
    console.log(`   内容：${el.text}`);
  });

  console.log('\n--- 表格信息 ---');
  if (pageAnalysis.tables.length > 0) {
    pageAnalysis.tables.forEach(table => {
      console.log(`表格 ${table.index + 1}:`);
      console.log(`  表头：${table.headers.join(' | ')}`);
      console.log(`  行数：${table.rowCount}`);
      if (table.firstRow.length > 0) {
        console.log(`  首行：${table.firstRow.join(' | ')}`);
      }
    });
  } else {
    console.log('未找到表格');
  }

  console.log('\n--- 输入框 ---');
  pageAnalysis.inputs.forEach(input => {
    console.log(`  - type="${input.type}" name="${input.name}" placeholder="${input.placeholder}"`);
  });

  console.log('\n--- 按钮 ---');
  pageAnalysis.buttons.slice(0, 20).forEach(btn => {
    console.log(`  - "${btn.text}" (type=${btn.type}, class=${btn.class})`);
  });

  console.log('\n--- 权限相关关键词 ---');
  pageAnalysis.permissionKeywords.forEach(kw => {
    console.log(`\n"${kw.keyword}" 相关:`);
    kw.matches.forEach(m => console.log(`  - ${m}`));
  });

  console.log('\n--- 页面文本内容（前 2000 字符）---');
  console.log(pageAnalysis.bodyText.substring(0, 2000));

  console.log('\n=== 完成 ===');
  await browser.close();
})().catch(err => {
  console.error('错误:', err);
  process.exit(1);
});
