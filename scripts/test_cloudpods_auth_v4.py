#!/usr/bin/env python3
"""
CloudPods 登录和 Dashboard 页面分析脚本 v4
使用 type 方法并简化流程
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

def main():
    all_api_calls = []

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

        # 监听 API 请求
        def on_request(req):
            url = req.url
            if 'api' in url.lower() or 'v1' in url.lower() or 'backend' in url.lower():
                all_api_calls.append({
                    'url': url,
                    'method': req.method,
                    'type': req.resource_type
                })

        page.on('request', on_request)

        print("\n" + "="*60)
        print("CloudPods 登录流程分析 v4")
        print("="*60)

        # 1. 访问登录页面
        print("\n>>> 步骤 1: 访问登录页面")
        page.goto('https://127.0.0.1/auth/login', wait_until='networkidle', timeout=30000)
        page.wait_for_timeout(2000)

        page.screenshot(path=os.path.join(OUTPUT_DIR, "login_page.png"))
        print("✓ 截图: login_page.png")

        # 2. 使用 keyboard 输入
        print("\n>>> 步骤 2: 填写登录信息")

        # 点击第一个输入框并输入用户名
        page.click('.ant-input >> nth=0')
        page.keyboard.type('admin', delay=50)
        print("✓ 输入用户名: admin")

        page.wait_for_timeout(500)

        # Tab 到密码框
        page.keyboard.press('Tab')
        page.keyboard.type('admin@123', delay=50)
        print("✓ 输入密码: admin@123")

        page.screenshot(path=os.path.join(OUTPUT_DIR, "login_filled.png"))
        print("✓ 截图: login_filled.png")

        # 3. 点击登录按钮
        print("\n>>> 步骤 3: 点击登录")
        page.click('.ant-btn-primary')
        print("✓ 点击登录按钮")

        page.wait_for_timeout(5000)

        current_url = page.url
        print(f"当前 URL: {current_url}")

        page.screenshot(path=os.path.join(OUTPUT_DIR, "after_login.png"))
        print("✓ 截图: after_login.png")

        if 'login' not in current_url.lower():
            print("✓ 登录成功!")

            # 4. 分析 Dashboard
            print("\n>>> 步骤 4: 分析 Dashboard")
            page.goto('https://127.0.0.1/dashboard', wait_until='networkidle', timeout=30000)
            page.wait_for_timeout(3000)

            page.screenshot(path=os.path.join(OUTPUT_DIR, "dashboard.png"), full_page=True)
            print("✓ 截图: dashboard.png")

            # 收集 Dashboard 数据
            dashboard = {
                "title": page.title(),
                "cards": [],
                "menu_items": []
            }

            # 卡片
            try:
                cards = page.locator('.ant-card').all()
                for card in cards[:10]:
                    try:
                        header = card.locator('.ant-card-head-title').inner_text()
                        dashboard["cards"].append(header)
                    except:
                        pass
            except:
                pass

            # 菜单
            try:
                menu_items = page.locator('.ant-menu-item').all()
                for m in menu_items:
                    try:
                        text = m.inner_text().strip()
                        if text:
                            dashboard["menu_items"].append(text)
                    except:
                        pass

                submenus = page.locator('.ant-menu-submenu-title').all()
                for sm in submenus:
                    try:
                        text = sm.inner_text().strip()
                        if text:
                            dashboard["menu_items"].append(f"[子菜单] {text}")
                    except:
                        pass
            except:
                pass

            save_json(dashboard, "dashboard.json")

            # 5. 展开菜单
            print("\n>>> 步骤 5: 展开菜单")
            try:
                submenus = page.locator('.ant-menu-submenu-title').all()
                for sm in submenus[:15]:
                    try:
                        sm.click()
                        page.wait_for_timeout(200)
                    except:
                        pass
            except:
                pass

            page.screenshot(path=os.path.join(OUTPUT_DIR, "menu_full.png"), full_page=True)
            print("✓ 截图: menu_full.png")

            # 收集完整菜单
            full_menu = []
            try:
                items = page.locator('.ant-menu-item, .ant-menu-submenu-title').all()
                for item in items:
                    try:
                        text = item.inner_text().strip()
                        if text:
                            full_menu.append(text)
                    except:
                        pass
            except:
                pass
            save_json(full_menu, "full_menu.json")

        # 保存 API 调用
        save_json(all_api_calls, "api_calls.json")

        # 提取 endpoints
        endpoints = set()
        for call in all_api_calls:
            url = call['url']
            match = re.search(r'https://[^/]+(/[^?]+)', url)
            if match:
                endpoints.add(match.group(1))
        save_json(sorted(list(endpoints)), "api_endpoints.json")

        print("\n" + "="*60)
        print("分析完成！输出目录:", OUTPUT_DIR)
        print("="*60)

        page.wait_for_timeout(3000)
        browser.close()

if __name__ == '__main__':
    main()