#!/usr/bin/env python3
"""
CloudPods 登录和 Dashboard 页面分析脚本 v3
使用最简单有效的选择器
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
        }) if ('api' in req.url.lower() or 'v1' in req.url.lower() or 'backend' in req.url.lower()) else None)

        page.on('response', lambda res: all_responses.append({
            'url': res.url,
            'status': res.status
        }) if ('api' in res.url.lower() or 'v1' in res.url.lower() or 'backend' in res.url.lower()) else None)

        print("\n" + "="*60)
        print("CloudPods 登录流程分析 v3")
        print("="*60)

        # 1. 直接访问带参数的登录页面
        print("\n>>> 步骤 1: 访问登录页面并执行登录")
        page.goto('https://127.0.0.1/auth/login?username=admin&fd_domain=Default', wait_until='networkidle')
        page.wait_for_timeout(2000)

        # 截图登录前
        page.screenshot(path=os.path.join(OUTPUT_DIR, "before_login.png"))
        print("✓ 截图: before_login.png")

        # 使用最简单的选择器 - 通过 placeholder 定位
        print("查找输入框...")

        # 等待页面完全加载
        page.wait_for_selector('.ant-input', timeout=10000)

        # 使用 nth 定位
        username_input = page.locator('.ant-input').nth(0)
        password_input = page.locator('.ant-input').nth(1)

        # 点击输入框使其聚焦
        username_input.click()
        page.wait_for_timeout(100)

        # 输入用户名
        username_input.fill('admin')
        print("✓ 填写用户名: admin")

        # 输入密码
        password_input.click()
        page.wait_for_timeout(100)
        password_input.fill('admin@123')
        print("✓ 填写密码: admin@123")

        # 截图填写后
        page.screenshot(path=os.path.join(OUTPUT_DIR, "after_fill.png"))
        print("✓ 截图: after_fill.png")

        # 点击登录按钮
        login_btn = page.locator('.ant-btn-primary').first
        login_btn.click()
        print("✓ 点击登录按钮")

        # 等待跳转
        page.wait_for_timeout(5000)

        current_url = page.url
        print(f"当前 URL: {current_url}")

        # 截图登录后
        page.screenshot(path=os.path.join(OUTPUT_DIR, "after_login.png"))
        print("✓ 截图: after_login.png")

        if 'login' not in current_url.lower():
            print("✓ 登录成功!")

            # 2. 分析 Dashboard
            print("\n>>> 步骤 2: 分析 Dashboard 页面")
            page.goto('https://127.0.0.1/dashboard', wait_until='networkidle')
            page.wait_for_timeout(3000)

            page.screenshot(path=os.path.join(OUTPUT_DIR, "dashboard.png"), full_page=True)
            print("✓ 截图: dashboard.png")

            # 分析 Dashboard 元素
            dashboard_data = {
                "title": page.title(),
                "url": page.url,
                "cards": [],
                "menu_items": [],
                "quick_links": []
            }

            # 卡片
            cards = page.locator('.ant-card').all()
            for card in cards[:15]:
                try:
                    header = card.locator('.ant-card-head-title').inner_text() or ''
                    body_text = card.locator('.ant-card-body').inner_text()[:200] or ''
                    dashboard_data["cards"].append({
                        "header": header,
                        "preview": body_text
                    })
                except:
                    pass

            # 菜单项
            menu_items = page.locator('.ant-menu-item').all()
            for m in menu_items:
                try:
                    text = m.inner_text().strip()
                    if text:
                        dashboard_data["menu_items"].append(text)
                except:
                    pass

            # 子菜单
            submenus = page.locator('.ant-menu-submenu-title').all()
            for sm in submenus:
                try:
                    text = sm.inner_text().strip()
                    if text:
                        dashboard_data["menu_items"].append(f"[子菜单] {text}")
                except:
                    pass

            save_json(dashboard_data, "dashboard_analysis.json")

            # 3. 展开所有菜单
            print("\n>>> 步骤 3: 展开菜单查看完整结构")
            for sm in submenus[:20]:
                try:
                    sm.click()
                    page.wait_for_timeout(200)
                except:
                    pass

            page.screenshot(path=os.path.join(OUTPUT_DIR, "menu_expanded.png"), full_page=True)
            print("✓ 截图: menu_expanded.png")

            # 收集完整菜单
            full_menu = []
            all_items = page.locator('.ant-menu-item, .ant-menu-submenu-title').all()
            for item in all_items:
                try:
                    text = item.inner_text().strip()
                    if text:
                        full_menu.append(text)
                except:
                    pass
            save_json(full_menu, "full_menu.json")

        else:
            print("登录可能失败，仍在登录页面")

        # 保存 API 数据
        save_json(all_api_calls, "api_calls.json")
        save_json(all_responses, "api_responses.json")

        # 提取 API endpoint
        endpoints = set()
        for call in all_api_calls:
            url = call['url']
            match = re.search(r'https://[^/]+(/[^?]+)', url)
            if match:
                path = match.group(1)
                endpoints.add(path)

        save_json(sorted(list(endpoints)), "api_endpoints.json")

        print("\n" + "="*60)
        print("分析完成！")
        print("输出目录:", OUTPUT_DIR)
        print("="*60)

        # 保持浏览器打开一会以便查看
        page.wait_for_timeout(5000)
        browser.close()

if __name__ == '__main__':
    main()