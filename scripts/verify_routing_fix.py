#!/usr/bin/env python3
"""
简单验证路由修复 - 增加超时时间
"""

from playwright.sync_api import sync_playwright
import time

def verify():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        page.set_default_timeout(60000)  # 60秒超时

        try:
            print("Step 1: 访问登录页面")
            page.goto('http://localhost:3000/login', timeout=60000)
            page.wait_for_load_state('domcontentloaded')
            time.sleep(2)

            print("Step 2: 登录")
            page.fill('input[placeholder="请输入用户名"]', 'admin')
            page.fill('input[placeholder="请输入密码"]', 'admin@123')
            page.click('button', timeout=60000)
            time.sleep(3)

            print(f"当前 URL: {page.url}")

            # 直接导航到虚拟机页面
            print("\nStep 3: 导航到虚拟机页面")
            page.goto('http://localhost:3000/compute/vms', timeout=60000)
            page.wait_for_load_state('domcontentloaded')
            time.sleep(2)

            vms_title = page.locator('.page-header h2').first
            vms_title_text = vms_title.inner_text() if vms_title.count() > 0 else "未找到标题"
            print(f"虚拟机页面标题: {vms_title_text}")

            # 导航到主机模版页面
            print("\nStep 4: 导航到主机模版页面")
            page.goto('http://localhost:3000/compute/host-templates', timeout=60000)
            page.wait_for_load_state('domcontentloaded')
            time.sleep(2)

            templates_title = page.locator('.page-header h2').first
            templates_title_text = templates_title.inner_text() if templates_title.count() > 0 else "未找到标题"
            print(f"主机模版页面标题: {templates_title_text}")

            # 验证结果
            print("\n=== 验证结果 ===")
            if templates_title_text == "主机模版":
                print("修复成功! 页面标题正确显示为'主机模版'")
            elif templates_title_text == "虚拟机管理":
                print("问题仍存在: 页面标题仍显示'虚拟机管理'")
            else:
                print(f"未知状态: 页面标题为'{templates_title_text}'")

            page.screenshot(path='/tmp/routing_verify_result.png')

        except Exception as e:
            print(f"错误: {e}")
            try:
                page.screenshot(path='/tmp/routing_verify_error.png')
            except:
                pass

        finally:
            browser.close()

if __name__ == "__main__":
    verify()