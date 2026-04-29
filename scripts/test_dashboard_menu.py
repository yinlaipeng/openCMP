#!/usr/bin/env python3
"""
测试 Dashboard 页面是否有左侧菜单栏
"""

from playwright.sync_api import sync_playwright
import os

OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_output"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context(viewport={'width': 1280, 'height': 900})
        page = context.new_page()

        print(">>> 1. 打开登录页面")
        page.goto('http://localhost:3000/login', wait_until='networkidle')
        page.wait_for_timeout(2000)

        print(">>> 2. 填写登录信息")
        page.locator('input').nth(0).fill('admin')
        page.locator('input').nth(1).fill('admin@123')

        print(">>> 3. 点击登录按钮")
        page.locator('button.el-button--primary').first.click()

        print(">>> 4. 等待跳转到 Dashboard")
        page.wait_for_timeout(5000)

        print(f"当前 URL: {page.url}")

        # 截图
        page.screenshot(path=os.path.join(OUTPUT_DIR, "dashboard_with_menu.png"))

        # 检查是否有 Layout 组件
        layout = page.locator('.layout-container, .app-wrapper, .sidebar-container, .el-menu').all()
        print(f"\n>>> 5. 检查 Layout 元素:")
        print(f"  找到 {len(layout)} 个 Layout 相关元素")

        # 检查侧边栏菜单
        sidebar = page.locator('.sidebar, .sidebar-container, aside, .el-aside').all()
        print(f"  找到 {len(sidebar)} 个侧边栏元素")

        # 检查菜单项
        menu_items = page.locator('.el-menu-item').all()
        print(f"  找到 {len(menu_items)} 个菜单项")

        if len(menu_items) > 0:
            print("\n菜单项列表:")
            for item in menu_items[:10]:
                print(f"    - {item.inner_text()}")
        else:
            print("\n❌ 没有找到菜单项！")

        # 检查统计卡片
        stat_cards = page.locator('.stat-card').all()
        print(f"\n  找到 {len(stat_cards)} 个统计卡片")

        print("\n测试完成")
        page.wait_for_timeout(3000)
        browser.close()

if __name__ == '__main__':
    main()