#!/usr/bin/env python3
"""
Cloudpods 负载均衡和 CDN 页面分析脚本
分析 4 个页面的完整设计结构
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
    btn_elements = page.locator('.toolbar button, .page-header button, .ant-btn, button.ant-btn').all()
    for btn in btn_elements[:20]:
        text = btn.inner_text().strip()
        if text and len(text) < 30:
            class_attr = btn.get_attribute('class') or ''
            btn_type = 'primary' if 'primary' in class_attr else 'default'
            disabled = 'disabled' in class_attr or btn.get_attribute('disabled')
            buttons.append({'text': text, 'type': btn_type, 'disabled': disabled})
    return buttons

def analyze_filters(page):
    """分析搜索/筛选区域"""
    filters = []
    form_items = page.locator('.ant-form-item, .search-form .ant-form-item').all()
    for item in form_items[:15]:
        try:
            label_el = item.locator('.ant-form-item-label label').first
            if label_el.count() > 0:
                label = label_el.inner_text().strip()
                if label:
                    filters.append(label)
        except:
            pass
    return filters

def analyze_dropdown_menu(page, button_text):
    """点击按钮查看下拉菜单"""
    menu_items = []
    try:
        btn = page.locator(f'button:has-text("{button_text}")').first
        if btn.count() > 0:
            btn.click()
            page.wait_for_timeout(500)
            items = page.locator('.ant-dropdown-menu-item, .ant-dropdown-item').all()
            for item in items:
                text = item.inner_text().strip()
                if text:
                    menu_items.append(text)
            page.keyboard.press('Escape')
            page.wait_for_timeout(300)
    except:
        pass
    return menu_items

def analyze_page(page, url, name, index):
    """分析单个页面"""
    print(f"\n{'='*60}")
    print(f"[{index}] 分析: {name}")
    print(f"URL: {url}")
    print(f"{'='*60}")

    page.goto(url, wait_until='networkidle', timeout=30000)
    page.wait_for_timeout(2000)

    # 截图
    screenshot = f'/tmp/cloudpods_lb_cdn_{index}_{name.replace(" ", "_")}.png'
    page.screenshot(path=screenshot, full_page=True)
    print(f"截图: {screenshot}")

    # 分析各部分
    headers = analyze_table(page)
    print(f"表头列: {headers}")

    buttons = analyze_toolbar(page)
    print(f"工具栏按钮: {[b['text'] for b in buttons[:10]]}")

    filters = analyze_filters(page)
    print(f"搜索筛选: {filters}")

    # 分析批量操作下拉菜单
    batch_menu = []
    for btn in buttons:
        if '批量' in btn['text'] or 'Batch' in btn['text'] or 'batch' in btn['text'].lower():
            batch_menu = analyze_dropdown_menu(page, btn['text'])
            print(f"批量操作菜单: {batch_menu}")
            break

    # 分析操作列
    actions_text = ''
    try:
        actions_cell = page.locator('table tbody tr td:last-child').first
        if actions_cell.count() > 0:
            actions_text = actions_cell.inner_text().strip()
            print(f"操作列: {actions_text}")
    except:
        pass

    return {
        'page': name,
        'url': url,
        'index': index,
        'headers': headers,
        'buttons': buttons,
        'filters': filters,
        'batch_menu': batch_menu,
        'actions': actions_text,
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

        print("="*60)
        print("登录 Cloudpods...")
        print("="*60)

        # 登录
        page.goto('https://127.0.0.1/auth/login', wait_until='networkidle', timeout=30000)
        page.wait_for_timeout(1500)

        # Ant Design 登录表单
        page.locator('input[placeholder*="username"]').fill('admin')
        page.locator('input[placeholder*="password"]').fill('admin@123')
        page.locator('button:has-text("Sign In")').click()

        page.wait_for_timeout(3000)
        page.wait_for_load_state('networkidle')
        print(f"登录成功，当前URL: {page.url}")

        # 分析各页面
        pages_to_analyze = [
            ('https://127.0.0.1/lb', '负载均衡实例', 1),
            ('https://127.0.0.1/lbacl', '访问控制', 2),
            ('https://127.0.0.1/lbcert', '证书', 3),
            ('https://127.0.0.1/cdn', 'CDN域名', 4),
        ]

        for url, name, index in pages_to_analyze:
            try:
                result = analyze_page(page, url, name, index)
                results.append(result)
            except Exception as e:
                print(f"分析 {name} 失败: {e}")
                results.append({'page': name, 'url': url, 'index': index, 'error': str(e)})

        # 保存结果
        output_file = '/tmp/cloudpods_lb_cdn_analysis.json'
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(results, f, ensure_ascii=False, indent=2)

        print(f"\n{'='*60}")
        print(f"分析完成！结果保存: {output_file}")
        print(f"分析了 {len(results)} 个页面")
        print(f"{'='*60}")

        browser.close()

if __name__ == '__main__':
    main()