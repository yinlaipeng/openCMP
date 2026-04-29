#!/usr/bin/env python3
"""
简化版路由诊断 - 先检查页面状态
"""

from playwright.sync_api import sync_playwright
import time

def diagnose():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        try:
            # 1. 访问登录页面
            print("Step 1: 检查登录页面")
            page.goto('http://localhost:3000/login')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            # 截图
            page.screenshot(path='/tmp/login_page_check.png')
            print(f"登录页 URL: {page.url}")

            # 检查表单元素
            inputs = page.locator('input').all()
            print(f"输入框数量: {len(inputs)}")
            for i, inp in enumerate(inputs):
                try:
                    placeholder = inp.get_attribute('placeholder')
                    type_attr = inp.get_attribute('type')
                    print(f"  Input {i}: placeholder={placeholder}, type={type_attr}")
                except:
                    pass

            buttons = page.locator('button').all()
            print(f"按钮数量: {len(buttons)}")
            for i, btn in enumerate(buttons):
                try:
                    text = btn.inner_text()
                    print(f"  Button {i}: text={text}")
                except:
                    pass

            # 2. 尝试登录
            print("\nStep 2: 尝试登录")

            # 使用更通用的 selector
            username_input = page.locator('input[type="text"], input:not([type="password"])').first
            password_input = page.locator('input[type="password"]').first
            login_button = page.locator('button').first

            if username_input.count() > 0:
                username_input.fill('admin')
                print("填写用户名: admin")
            else:
                print("未找到用户名输入框")

            if password_input.count() > 0:
                password_input.fill('admin@123')
                print("填写密码: admin@123")
            else:
                print("未找到密码输入框")

            if login_button.count() > 0:
                login_button.click()
                print("点击登录按钮")
                time.sleep(3)
                page.wait_for_load_state('networkidle')
                print(f"登录后 URL: {page.url}")
                page.screenshot(path='/tmp/after_login.png')
            else:
                print("未找到登录按钮")

            # 3. 如果登录成功，测试路由切换
            if '/login' not in page.url:
                print("\nStep 3: 测试路由切换")

                # 直接导航到虚拟机页面
                page.goto('http://localhost:3000/compute/vms')
                page.wait_for_load_state('networkidle')
                time.sleep(2)

                vms_url = page.url
                vms_content = page.content()
                print(f"虚拟机页面 URL: {vms_url}")
                print(f"页面包含 '虚拟机': {'虚拟机' in vms_content}")

                page.screenshot(path='/tmp/vms_page.png')

                # 点击其他菜单 - 直接导航测试
                print("\nStep 4: 直接导航到主机模版页面")
                page.goto('http://localhost:3000/compute/host-templates')
                page.wait_for_load_state('networkidle')
                time.sleep(2)

                templates_url = page.url
                templates_content = page.content()
                print(f"主机模版 URL: {templates_url}")
                print(f"页面包含 '模版': {'模版' in templates_content}")
                print(f"页面包含 '虚拟机': {'虚拟机' in templates_content}")

                page.screenshot(path='/tmp/templates_page_direct.png')

                # 现在测试点击侧边栏菜单
                print("\nStep 5: 测试点击侧边栏菜单")

                # 先回到虚拟机页面
                page.goto('http://localhost:3000/compute/vms')
                page.wait_for_load_state('networkidle')
                time.sleep(2)

                # 寻找侧边栏菜单
                sidebar_items = page.locator('.el-menu-item').all()
                print(f"侧边栏菜单项数量: {len(sidebar_items)}")

                # 点击一个不同的菜单项
                for item in sidebar_items:
                    try:
                        text = item.inner_text()
                        if '镜像' in text or '安全组' in text:
                            print(f"点击菜单: {text}")
                            item.click()
                            time.sleep(3)
                            page.wait_for_load_state('networkidle')

                            new_url = page.url
                            new_content = page.content()
                            print(f"点击后 URL: {new_url}")
                            print(f"页面内容变化检测:")
                            print(f"  - 包含 '{text}': {text in new_content}")
                            print(f"  - 包含 '虚拟机': {'虚拟机' in new_content}")

                            page.screenshot(path=f'/tmp/after_click_{text.replace("/", "_")}.png')
                            break
                    except Exception as e:
                        print(f"菜单点击错误: {e}")

                print("\n=== 测试完成 ===")

        except Exception as e:
            print(f"错误: {e}")
            page.screenshot(path='/tmp/error_screenshot.png')

        finally:
            browser.close()

if __name__ == "__main__":
    diagnose()