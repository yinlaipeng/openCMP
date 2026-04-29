#!/usr/bin/env python3
"""
openCMP 登录调试脚本 - 实时检查 localStorage
"""

from playwright.sync_api import sync_playwright
import os
import time
import json

OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_output/opencmp_auth_dashboard"

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context(viewport={'width': 1280, 'height': 900})
        page = context.new_page()

        # 监听控制台
        console_messages = []
        page.on('console', lambda msg: console_messages.append({
            'type': msg.type,
            'text': msg.text
        }))

        # 监听 localStorage 变化
        def check_storage():
            return page.evaluate("() => { return { token: localStorage.getItem('token'), user: localStorage.getItem('user') }; }")

        print(">>> 打开登录页面")
        page.goto('http://localhost:3000/login', wait_until='networkidle')
        page.wait_for_timeout(2000)

        # 初始 localStorage 检查
        storage = check_storage()
        print(f"初始 localStorage: token={storage['token']}")

        print(">>> 填写登录信息")
        page.locator('input').nth(0).fill('admin')
        page.locator('input').nth(1).fill('admin@123')

        print(">>> 点击登录按钮")
        page.locator('button.el-button--primary').first.click()

        # 等待一小段时间后检查
        page.wait_for_timeout(500)
        storage = check_storage()
        print(f"点击后 500ms: token={storage['token'][:50] if storage['token'] else 'None'}")

        page.wait_for_timeout(1000)
        storage = check_storage()
        print(f"点击后 1s: token={storage['token'][:50] if storage['token'] else 'None'}")

        page.wait_for_timeout(2000)
        storage = check_storage()
        print(f"点击后 2s: token={storage['token'][:50] if storage['token'] else 'None'}")

        page.wait_for_timeout(5000)
        storage = check_storage()
        print(f"点击后 5s: token={storage['token'][:50] if storage['token'] else 'None'}")

        print(f"\n最终 URL: {page.url}")
        print(f"页面标题: {page.title()}")

        # 检查控制台消息
        print(f"\n控制台消息:")
        for msg in console_messages:
            if msg['type'] in ['log', 'error']:
                print(f"  [{msg['type']}] {msg['text'][:100]}")

        # 截图
        page.screenshot(path=os.path.join(OUTPUT_DIR, "debug_storage.png"))

        page.wait_for_timeout(3000)
        browser.close()

if __name__ == '__main__':
    main()