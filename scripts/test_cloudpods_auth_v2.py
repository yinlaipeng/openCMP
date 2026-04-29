#!/usr/bin/env python3
"""
CloudPods 登录和 Dashboard 页面分析脚本 v2
使用正确的 Ant Design 选择器
"""

from playwright.sync_api import sync_playwright
import json
import os
import re

OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_output/cloudpods_auth"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def save_json(data, filename):
    path = os.path.join(OUTPUT_DIR, filename)
    with open(path, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"✓ 保存: {path}")

def analyze_page(page, name):
    """分析页面结构和 API"""
    print(f"\n=== 分析页面: {name} ===")

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{name}_screenshot.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"✓ 截图: {screenshot_path}")

    # 页面标题
    title = page.title()
    print(f"页面标题: {title}")

    # 分析元素
    elements = {
        "inputs": [],
        "buttons": [],
        "forms": [],
        "titles": [],
        "menu_items": [],
        "cards": []
    }

    # 输入框 - Ant Design
    inputs = page.locator('.ant-input, input').all()
    for inp in inputs:
        try:
            elements["inputs"].append({
                "type": inp.get_attribute('type') or 'text',
                "placeholder": inp.get_attribute('placeholder') or '',
                "class": inp.get_attribute('class') or ''
            })
        except:
            pass

    # 按钮 - Ant Design
    buttons = page.locator('.ant-btn, button').all()
    for btn in buttons:
        try:
            elements["buttons"].append({
                "text": btn.inner_text().strip(),
                "type": btn.get_attribute('type') or '',
                "class": btn.get_attribute('class') or ''
            })
        except:
            pass

    # 标题
    titles = page.locator('h1, h2, h3, .ant-page-header-title, .title').all()
    for t in titles:
        try:
            text = t.inner_text().strip()
            if text:
                elements["titles"].append(text)
        except:
            pass

    # 卡片
    cards = page.locator('.ant-card, .card, .widget').all()
    for card in cards[:10]:
        try:
            header = card.locator('.ant-card-head-title').inner_text() or ''
            elements["cards"].append({"header": header})
        except:
            pass

    # 菜单
    menus = page.locator('.ant-menu-item, .ant-menu-submenu-title, nav a').all()
    for m in menus[:30]:
        try:
            text = m.inner_text().strip()
            if text and len(text) < 50:
                elements["menu_items"].append(text)
        except:
            pass

    save_json(elements, f"{name}_elements.json")
    return elements

def main():
    all_api_calls = []
    all_responses = []

    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=False,
            args=['--ignore-certificate-errors', '--disable-web-security']
        )

        context = browser.new_context(
            ignore_https_errors=True,
            viewport={'width': 1280, 'height': 900}
        )

        page = context.new_page()

        # 监听所有 API 请求
        page.on('request', lambda req: all_api_calls.append({
            'url': req.url,
            'method': req.method,
            'type': req.resource_type
        }) if ('api' in req.url.lower() or 'auth' in req.url.lower() or 'v1' in req.url.lower()) else None)

        page.on('response', lambda res: all_responses.append({
            'url': res.url,
            'status': res.status
        }) if ('api' in res.url.lower() or 'auth' in res.url.lower() or 'v1' in res.url.lower()) else None)

        print("\n" + "="*60)
        print("CloudPods 登录流程分析 v2")
        print("="*60)

        # 1. 访问登录页面
        print("\n>>> 步骤 1: 访问登录页面")
        page.goto('https://127.0.0.1/auth/login', wait_until='networkidle')
        page.wait_for_timeout(1000)
        analyze_page(page, 'login_page')

        # 2. 访问账号选择页面
        print("\n>>> 步骤 2: 访问账号选择页面")
        page.goto('https://127.0.0.1/auth/login/chooser', wait_until='networkidle')
        page.wait_for_timeout(1000)
        analyze_page(page, 'chooser_page')

        # 3. 回到登录页面执行登录
        print("\n>>> 步骤 3: 执行登录操作")
        page.goto('https://127.0.0.1/auth/login?username=admin&fd_domain=Default', wait_until='networkidle')
        page.wait_for_timeout(1000)

        # 使用 Ant Design 选择器
        username_input = page.locator('input.ant-input[placeholder*="用户名"]').first
        password_input = page.locator('input.ant-input[type="password"]').first

        if username_input:
            username_input.fill('admin')
            print("✓ 填写用户名: admin")

        if password_input:
            password_input.fill('admin@123')
            print("✓ 填写密码: admin@123")

        # 点击登录按钮
        login_btn = page.locator('button.ant-btn-primary:has-text("登录")').first
        if login_btn:
            login_btn.click()
            print("✓ 点击登录按钮")
            page.wait_for_timeout(3000)

            current_url = page.url
            print(f"当前 URL: {current_url}")

            # 检查是否成功
            if 'dashboard' in current_url.lower() or 'home' in current_url.lower() or current_url != 'https://127.0.0.1/auth/login':
                print("✓ 登录成功!")

                # 4. 分析 Dashboard
                print("\n>>> 步骤 4: 分析 Dashboard 页面")
                page.goto('https://127.0.0.1/dashboard', wait_until='networkidle')
                page.wait_for_timeout(2000)
                analyze_page(page, 'dashboard')

                # 5. 分析侧边菜单结构
                print("\n>>> 步骤 5: 分析完整菜单")
                # 展开所有菜单
                submenus = page.locator('.ant-menu-submenu-title').all()
                for sm in submenus[:10]:
                    try:
                        sm.click()
                        page.wait_for_timeout(300)
                    except:
                        pass

                analyze_page(page, 'dashboard_menu_expanded')

                # 保存完整菜单
                menu_data = []
                all_menu_items = page.locator('.ant-menu-item, .ant-menu-submenu-title').all()
                for item in all_menu_items:
                    try:
                        text = item.inner_text().strip()
                        if text:
                            menu_data.append(text)
                    except:
                        pass
                save_json(menu_data, "full_menu_structure.json")

        # 保存 API 调用
        save_json(all_api_calls, "all_api_requests.json")
        save_json(all_responses, "all_api_responses.json")

        # 从 URL 中提取 API endpoint
        api_endpoints = []
        for call in all_api_calls:
            url = call['url']
            # 提取路径
            if '/api/' in url or '/v1/' in url:
                match = re.search(r'(https://[^/]+)?(/[^?]+)', url)
                if match:
                    path = match.group(2)
                    api_endpoints.append({
                        'path': path,
                        'method': call['method'],
                        'type': call['type']
                    })

        save_json(api_endpoints, "api_endpoints_summary.json")

        print("\n" + "="*60)
        print("分析完成！输出目录:")
        print(OUTPUT_DIR)
        print("="*60)

        browser.close()

if __name__ == '__main__':
    main()