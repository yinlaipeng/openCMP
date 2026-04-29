#!/usr/bin/env python3
"""
Cloudpods 页面验证脚本 - 最终检查
验证安全组、IP子网、弹性公网IP、密钥页面设计
"""

import json
from playwright.sync_api import sync_playwright

def analyze_table(page):
    """分析表格结构"""
    headers = []
    header_cells = page.locator('table thead tr th, .ant-table-thead th').all()
    for cell in header_cells:
        text = cell.inner_text().strip()
        if text and text not in ['操作', 'Operations', '']:
            headers.append(text)
    return headers

def analyze_toolbar(page):
    """分析工具栏按钮"""
    buttons = []
    btn_elements = page.locator('.toolbar button, .page-header button, .ant-btn').all()
    for btn in btn_elements[:15]:
        text = btn.inner_text().strip()
        if text and len(text) < 25:
            class_attr = btn.get_attribute('class') or ''
            btn_type = 'primary' if 'primary' in class_attr else 'default'
            buttons.append({'text': text, 'type': btn_type})
    return buttons

def analyze_filters(page):
    """分析搜索/筛选区域"""
    filters = []
    # 查找筛选表单
    form_items = page.locator('.filter-card .el-form-item, .ant-form-item').all()
    for item in form_items[:10]:
        label = item.locator('label, .el-form-item__label, .ant-form-item-label').first.inner_text().strip()
        if label:
            filters.append(label)
    return filters

def analyze_page(page, url, name):
    """分析单个页面"""
    print(f"\n{'='*50}")
    print(f"分析: {name}")
    print(f"URL: {url}")
    print(f"{'='*50}")

    page.goto(url, wait_until='networkidle')
    page.wait_for_timeout(2000)

    # 截图
    screenshot = f'/tmp/cloudpods_final_{name}.png'
    page.screenshot(path=screenshot, full_page=True)
    print(f"截图: {screenshot}")

    # 分析各部分
    headers = analyze_table(page)
    print(f"表头: {headers}")

    buttons = analyze_toolbar(page)
    print(f"按钮: {[b['text'] for b in buttons[:8]]}")

    filters = analyze_filters(page)
    print(f"筛选: {filters}")

    # 操作列
    actions = page.locator('table tbody tr td:last-child').first.inner_text().strip()
    print(f"操作列: {actions}")

    return {
        'page': name,
        'url': url,
        'headers': headers,
        'buttons': buttons,
        'filters': filters,
        'actions': actions,
        'screenshot': screenshot
    }

def main():
    results = []

    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--disable-web-security']
        )
        context = browser.new_context(ignore_https_errors=True)
        page = context.new_page()

        print("登录 Cloudpods...")

        # 登录
        page.goto('https://127.0.0.1/auth/login', wait_until='networkidle')
        page.wait_for_timeout(1500)

        # Ant Design 登录表单
        page.locator('input[placeholder*="username"]').fill('admin')
        page.locator('input[placeholder*="password"]').fill('admin@123')
        page.locator('button:has-text("Sign In")').click()

        page.wait_for_timeout(3000)
        page.wait_for_load_state('networkidle')
        print(f"登录成功，当前URL: {page.url}")

        # 分析各页面
        pages = [
            ('https://127.0.0.1/secgroup', '安全组'),
            ('https://127.0.0.1/network2', 'IP子网'),
            ('https://127.0.0.1/eip2', '弹性公网IP'),
            ('https://127.0.0.1/keypair', '密钥'),
        ]

        for url, name in pages:
            try:
                result = analyze_page(page, url, name)
                results.append(result)
            except Exception as e:
                print(f"分析 {name} 失败: {e}")
                results.append({'page': name, 'error': str(e)})

        # 保存结果
        with open('/tmp/cloudpods_final_analysis.json', 'w', encoding='utf-8') as f:
            json.dump(results, f, ensure_ascii=False, indent=2)

        print(f"\n结果保存: /tmp/cloudpods_final_analysis.json")
        browser.close()

if __name__ == '__main__':
    main()