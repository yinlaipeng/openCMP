#!/usr/bin/env python3
"""
验证侧边栏菜单点击切换路由
"""

from playwright.sync_api import sync_playwright
import time

def verify_sidebar_click():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        page.set_default_timeout(60000)

        results = []

        try:
            print("Step 1: 登录")
            page.goto('http://localhost:3000/login', timeout=60000)
            page.wait_for_load_state('domcontentloaded')
            time.sleep(2)

            page.fill('input[placeholder="请输入用户名"]', 'admin')
            page.fill('input[placeholder="请输入密码"]', 'admin@123')
            page.click('button', timeout=60000)
            time.sleep(3)

            # 2. 先导航到虚拟机页面
            print("\nStep 2: 导航到虚拟机页面")
            page.goto('http://localhost:3000/compute/vms', timeout=60000)
            page.wait_for_load_state('domcontentloaded')
            time.sleep(3)

            vms_title = page.locator('.page-header h2').first.inner_text()
            print(f"虚拟机页面标题: {vms_title}")
            results.append(("虚拟机页面初始", vms_title, page.url))

            # 3. 点击侧边栏的"主机模版"菜单项
            print("\nStep 3: 点击侧边栏菜单切换")

            # 展开主机子菜单（如果需要）
            compute_host_submenu = page.locator('.el-sub-menu:has-text("主机")')
            if compute_host_submenu.count() > 0:
                # 检查子菜单是否已展开
                submenu_class = compute_host_submenu.first.get_attribute('class')
                if 'is-opened' not in (submenu_class or ''):
                    print("展开主机子菜单")
                    compute_host_submenu.first.click()
                    time.sleep(1)

            # 点击主机模版菜单项
            host_templates_item = page.locator('.el-menu-item:has-text("主机模版")')
            if host_templates_item.count() > 0:
                print("点击主机模版菜单项")
                host_templates_item.first.click()
                time.sleep(3)
                page.wait_for_load_state('domcontentloaded')
                time.sleep(2)

                new_title = page.locator('.page-header h2').first.inner_text()
                new_url = page.url
                print(f"点击后页面标题: {new_title}")
                print(f"点击后 URL: {new_url}")
                results.append(("点击菜单后", new_title, new_url))
            else:
                print("未找到主机模版菜单项")
                results.append(("点击菜单后", "未找到菜单项", page.url))

            # 4. 点击镜像菜单
            print("\nStep 4: 点击镜像菜单")
            images_item = page.locator('.el-menu-item:has-text("系统镜像")')
            if images_item.count() > 0:
                print("点击系统镜像菜单项")
                images_item.first.click()
                time.sleep(3)
                page.wait_for_load_state('domcontentloaded')
                time.sleep(2)

                images_title = page.locator('.page-header h2').first.inner_text()
                images_url = page.url
                print(f"镜像页面标题: {images_title}")
                print(f"镜像页面 URL: {images_url}")
                results.append(("点击镜像后", images_title, images_url))
            else:
                print("未找到镜像菜单项")

            # 5. 验证结果
            print("\n=== 验证结果 ===")
            success = True
            for step, title, url in results:
                print(f"{step}: 标题={title}, URL={url}")
                if step == "点击菜单后" and title == "虚拟机管理":
                    success = False
                    print("问题仍存在: 点击菜单后标题仍为虚拟机管理")
                if step == "点击镜像后" and title == "虚拟机管理":
                    success = False
                    print("问题仍存在: 点击镜像后标题仍为虚拟机管理")

            if success:
                print("\n修复完全成功! 侧边栏菜单点击切换正常工作")

            page.screenshot(path='/tmp/sidebar_click_verify.png')

        except Exception as e:
            print(f"错误: {e}")
            try:
                page.screenshot(path='/tmp/sidebar_click_error.png')
            except:
                pass

        finally:
            browser.close()

        return results

if __name__ == "__main__":
    verify_sidebar_click()