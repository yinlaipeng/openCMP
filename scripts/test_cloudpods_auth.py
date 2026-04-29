#!/usr/bin/env python3
"""
CloudPods 登录和 Dashboard 页面分析脚本
分析页面设计、API、组件结构，用于指导 openCMP 开发
"""

from playwright.sync_api import sync_playwright
import json
import os

# 输出目录
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_output/cloudpods_auth"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def save_json(data, filename):
    path = os.path.join(OUTPUT_DIR, filename)
    with open(path, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"✓ 保存: {path}")

def analyze_page(page, url, name):
    """分析页面结构和 API"""
    print(f"\n=== 分析页面: {name} ===")
    print(f"URL: {url}")

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{name}_screenshot.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"✓ 截图: {screenshot_path}")

    # 页面内容
    html_path = os.path.join(OUTPUT_DIR, f"{name}_html.html")
    with open(html_path, 'w', encoding='utf-8') as f:
        f.write(page.content())
    print(f"✓ HTML: {html_path}")

    # 分析元素
    elements = {
        "inputs": [],
        "buttons": [],
        "links": [],
        "forms": [],
        "titles": [],
        "api_calls": []
    }

    # 输入框
    inputs = page.locator('input').all()
    for inp in inputs:
        try:
            elements["inputs"].append({
                "type": inp.get_attribute('type') or 'text',
                "name": inp.get_attribute('name') or '',
                "placeholder": inp.get_attribute('placeholder') or '',
                "id": inp.get_attribute('id') or '',
                "class": inp.get_attribute('class') or ''
            })
        except:
            pass

    # 按钮
    buttons = page.locator('button').all()
    for btn in buttons:
        try:
            elements["buttons"].append({
                "text": btn.inner_text() or '',
                "type": btn.get_attribute('type') or '',
                "class": btn.get_attribute('class') or ''
            })
        except:
            pass

    # 链接
    links = page.locator('a').all()
    for link in links[:20]:  # 只取前20个
        try:
            elements["links"].append({
                "text": link.inner_text() or '',
                "href": link.get_attribute('href') or ''
            })
        except:
            pass

    # 表单
    forms = page.locator('form').all()
    for form in forms:
        try:
            elements["forms"].append({
                "action": form.get_attribute('action') or '',
                "method": form.get_attribute('method') or '',
                "id": form.get_attribute('id') or ''
            })
        except:
            pass

    # 标题
    titles = page.locator('h1, h2, h3, .title, .el-page-header__title, .el-card__header').all()
    for t in titles:
        try:
            text = t.inner_text()
            if text:
                elements["titles"].append(text)
        except:
            pass

    save_json(elements, f"{name}_elements.json")
    return elements

def main():
    with sync_playwright() as p:
        # 启动浏览器，忽略 SSL 证书
        browser = p.chromium.launch(
            headless=False,  # 可见模式便于调试
            args=['--ignore-certificate-errors', '--disable-web-security']
        )

        context = browser.new_context(
            ignore_https_errors=True,
            viewport={'width': 1280, 'height': 900}
        )

        page = context.new_page()

        # 监听 API 请求
        api_calls = []
        page.on('request', lambda req: api_calls.append({
            'url': req.url,
            'method': req.method,
            'resource_type': req.resource_type
        }) if '/api/' in req.url or '/auth/' in req.url else None)

        # 监听响应
        responses = []
        page.on('response', lambda res: responses.append({
            'url': res.url,
            'status': res.status,
            'body': res.text() if res.status < 400 else ''
        }) if '/api/' in res.url or '/auth/' in res.url else None)

        print("\n" + "="*60)
        print("CloudPods 登录流程分析")
        print("="*60)

        # 1. 访问登录页面
        print("\n>>> 步骤 1: 访问登录页面")
        page.goto('https://127.0.0.1/auth/login', wait_until='networkidle')
        page.wait_for_timeout(2000)
        analyze_page(page, 'https://127.0.0.1/auth/login', 'login_page')

        # 2. 查看账号选择器页面
        print("\n>>> 步骤 2: 访问账号选择页面")
        page.goto('https://127.0.0.1/auth/login/chooser', wait_until='networkidle')
        page.wait_for_timeout(2000)
        analyze_page(page, 'https://127.0.0.1/auth/login/chooser', 'chooser_page')

        # 3. 带参数的登录页面
        print("\n>>> 步骤 3: 访问带参数登录页面")
        page.goto('https://127.0.0.1/auth/login?username=admin&fd_domain=Default', wait_until='networkidle')
        page.wait_for_timeout(2000)
        analyze_page(page, 'https://127.0.0.1/auth/login?username=admin&fd_domain=Default', 'login_with_params')

        # 4. 执行登录
        print("\n>>> 步骤 4: 执行登录操作")

        # 查找并填写用户名
        username_input = page.locator('input[type="text"], input[name="username"], input[placeholder*="用户"]').first
        if username_input:
            username_input.fill('admin')
            print("✓ 填写用户名: admin")

        # 查找并填写密码
        password_input = page.locator('input[type="password"]').first
        if password_input:
            password_input.fill('admin@123')
            print("✓ 填写密码: admin@123")

        # 查找登录按钮
        login_btn = page.locator('button:has-text("登录"), button:has-text("Login"), button[type="submit"]').first
        if login_btn:
            login_btn.click()
            print("✓ 点击登录按钮")

            # 等待跳转
            page.wait_for_timeout(3000)

            # 检查是否登录成功
            current_url = page.url
            print(f"当前 URL: {current_url}")

            if 'dashboard' in current_url or 'home' in current_url:
                print("✓ 登录成功!")

                # 5. 分析 Dashboard
                print("\n>>> 步骤 5: 分析 Dashboard 页面")
                analyze_page(page, current_url, 'dashboard')

                # 6. 分析菜单结构
                print("\n>>> 步骤 6: 分析菜单结构")
                menu_items = page.locator('.el-menu-item, .el-submenu__title, nav a, .sidebar a').all()
                menu_data = []
                for item in menu_items:
                    try:
                        menu_data.append({
                            "text": item.inner_text(),
                            "href": item.get_attribute('href') or ''
                        })
                    except:
                        pass
                save_json(menu_data, "menu_structure.json")

                # 7. 分析 Dashboard 卡片/组件
                print("\n>>> 步骤 7: 分析 Dashboard 卡片")
                cards = page.locator('.el-card, .card, .dashboard-card, .widget').all()
                card_data = []
                for card in cards:
                    try:
                        card_data.append({
                            "header": card.locator('.el-card__header, .card-header').inner_text() or '',
                            "content": card.inner_text()[:500]
                        })
                    except:
                        pass
                save_json(card_data, "dashboard_cards.json")

        # 保存 API 调用记录
        save_json(api_calls, "api_requests.json")
        save_json(responses, "api_responses.json")

        print("\n" + "="*60)
        print("分析完成！输出文件保存在:")
        print(OUTPUT_DIR)
        print("="*60)

        browser.close()

if __name__ == '__main__':
    main()