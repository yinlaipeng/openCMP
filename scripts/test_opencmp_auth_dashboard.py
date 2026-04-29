#!/usr/bin/env python3
"""
openCMP 登录和 Dashboard 测试脚本
测试登录流程和 Dashboard 页面功能
"""

from playwright.sync_api import sync_playwright
import json
import os
import time

OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_output/opencmp_auth_dashboard"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def save_json(data, filename):
    path = os.path.join(OUTPUT_DIR, filename)
    with open(path, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"✓ 保存: {path}")

def main():
    test_results = {
        "login_page": {},
        "chooser_page": {},
        "dashboard_page": {},
        "auth_api": {},
        "summary": {}
    }

    api_errors = []
    console_errors = []

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context(viewport={'width': 1280, 'height': 900})
        page = context.new_page()

        # 监听 API 请求和响应
        api_calls = []
        page.on('response', lambda res: api_calls.append({
            'url': res.url,
            'status': res.status
        }) if '/api/' in res.url or '/auth/' in res.url else None)

        # 监听控制台错误
        page.on('console', lambda msg: console_errors.append(msg.text) if msg.type == 'error' else None)

        print("\n" + "="*60)
        print("openCMP 登录和 Dashboard 测试")
        print("="*60)

        # 1. 测试登录页面
        print("\n>>> 测试 1: 登录页面")
        try:
            page.goto('http://localhost:3000/login', wait_until='networkidle')
            page.wait_for_timeout(2000)

            # 截图
            page.screenshot(path=os.path.join(OUTPUT_DIR, "login_page.png"))
            print("✓ 登录页面加载成功")

            # 检查元素
            inputs = page.locator('input').all()
            buttons = page.locator('button').all()

            test_results["login_page"] = {
                "status": "success",
                "input_count": len(inputs),
                "button_count": len(buttons),
                "title": page.title()
            }
        except Exception as e:
            test_results["login_page"] = {"status": "error", "message": str(e)}
            print(f"✗ 登录页面加载失败: {e}")

        # 2. 测试账号选择器页面
        print("\n>>> 测试 2: 账号选择器页面")
        try:
            page.goto('http://localhost:3000/auth/login/chooser', wait_until='networkidle')
            page.wait_for_timeout(2000)

            page.screenshot(path=os.path.join(OUTPUT_DIR, "chooser_page.png"))
            print("✓ 账号选择器页面加载成功")

            test_results["chooser_page"] = {
                "status": "success",
                "title": page.title()
            }
        except Exception as e:
            test_results["chooser_page"] = {"status": "error", "message": str(e)}
            print(f"✗ 账号选择器页面加载失败: {e}")

        # 3. 测试登录流程
        print("\n>>> 测试 3: 登录流程")
        try:
            page.goto('http://localhost:3000/login', wait_until='networkidle')
            page.wait_for_timeout(1000)

            # 填写登录信息
            page.locator('input').nth(0).fill('admin')
            page.locator('input').nth(1).fill('admin@123')
            print("✓ 填写登录信息")

            page.screenshot(path=os.path.join(OUTPUT_DIR, "login_filled.png"))

            # 点击登录按钮
            page.locator('button.el-button--primary').first.click()
            print("✓ 点击登录按钮")

            page.wait_for_timeout(5000)

            current_url = page.url
            print(f"当前 URL: {current_url}")

            if 'dashboard' in current_url:
                print("✓ 登录成功，跳转到 Dashboard")

                # 截图
                page.screenshot(path=os.path.join(OUTPUT_DIR, "dashboard.png"))

                # 检查 Dashboard 元素
                cards = page.locator('.el-card').all()
                stat_cards = page.locator('.stat-card').all()

                test_results["dashboard_page"] = {
                    "status": "success",
                    "url": current_url,
                    "card_count": len(cards),
                    "stat_card_count": len(stat_cards),
                    "title": page.title()
                }
            else:
                test_results["dashboard_page"] = {"status": "error", "url": current_url}
                print(f"✗ 登录后未跳转到 Dashboard: {current_url}")

        except Exception as e:
            test_results["login_flow"] = {"status": "error", "message": str(e)}
            print(f"✗ 登录流程失败: {e}")

        # 4. 检查 API 调用
        print("\n>>> 测试 4: API 调用检查")
        auth_apis = ['/auth/login', '/auth/user', '/auth/permissions', '/auth/stats']
        api_status = {}

        for api_call in api_calls:
            url = api_call['url']
            status = api_call['status']
            for api in auth_apis:
                if api in url:
                    api_status[api] = status
                    if status >= 400:
                        api_errors.append(f"{url}: {status}")

        test_results["auth_api"] = api_status

        for api, status in api_status.items():
            print(f"  {api}: {status}")

        # 5. 检查控制台错误
        print("\n>>> 测试 5: 控制台错误检查")
        if console_errors:
            print(f"✗ 发现 {len(console_errors)} 个控制台错误")
            for err in console_errors[:5]:
                print(f"  - {err[:100]}")
        else:
            print("✓ 无控制台错误")

        # 统计结果
        success_count = sum(1 for k, v in test_results.items() if isinstance(v, dict) and v.get('status') == 'success')
        total_count = len([k for k in test_results if k != 'summary'])

        test_results["summary"] = {
            "success_count": success_count,
            "total_count": total_count,
            "api_errors": api_errors,
            "console_errors": console_errors,
            "pass_rate": f"{success_count}/{total_count}"
        }

        save_json(test_results, "test_results.json")

        print("\n" + "="*60)
        print(f"测试完成: {success_count}/{total_count} 成功")
        print("输出目录:", OUTPUT_DIR)
        print("="*60)

        page.wait_for_timeout(3000)
        browser.close()

if __name__ == '__main__':
    main()