#!/usr/bin/env python3
"""
openCMP 云账号同步功能详细测试 v2
正确模拟用户同步操作流程
"""

from playwright.sync_api import sync_playwright
import json
import os
import time

OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_output/cloud_account_sync"
os.makedirs(OUTPUT_DIR, exist_ok=True)

def save_json(data, filename):
    path = os.path.join(OUTPUT_DIR, filename)
    with open(path, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"✓ 保存: {path}")

def main():
    test_results = {}
    api_calls = []
    console_errors = []

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context(viewport={'width': 1400, 'height': 900})
        page = context.new_page()

        # 监听 API 请求
        page.on('request', lambda req: api_calls.append({
            'url': req.url,
            'method': req.method
        }) if '/api/' in req.url and 'cloud-account' in req.url else None)

        # 监听控制台错误
        page.on('console', lambda msg: console_errors.append(msg.text) if msg.type == 'error' else None)

        print("\n" + "="*70)
        print("openCMP 云账号同步功能详细测试 v2")
        print("="*70)

        # ========== 阶段1: 登录 ==========
        print("\n>>> 阶段1: 登录")
        page.goto('http://localhost:3000/login', wait_until='networkidle')
        page.wait_for_timeout(2000)
        page.locator('input').nth(0).fill('admin')
        page.locator('input').nth(1).fill('admin@123')
        page.locator('button.el-button--primary').first.click()
        page.wait_for_timeout(5000)
        print(f"当前 URL: {page.url}")
        page.screenshot(path=os.path.join(OUTPUT_DIR, "login.png"))

        # ========== 阶段2: 进入云账号页面 ==========
        print("\n>>> 阶段2: 云账号页面")
        page.goto('http://localhost:3000/cloud-accounts', wait_until='networkidle')
        page.wait_for_timeout(3000)
        page.screenshot(path=os.path.join(OUTPUT_DIR, "cloud_accounts.png"))

        # ========== 阶段3: 检查云账号列表 ==========
        print("\n>>> 阶段3: 检查云账号列表")

        # 获取表格数据
        table_rows = page.locator('.el-table__row').all()
        print(f"表格行数: {len(table_rows)}")

        accounts_info = []
        for row in table_rows:
            cells = row.locator('.el-table__cell').all()
            row_data = {}
            for i, cell in enumerate(cells):
                text = cell.inner_text().strip()
                if i == 1: row_data['id'] = text
                elif i == 2: row_data['name'] = text
                elif i == 3: row_data['status'] = text
                elif i == 4: row_data['enabled'] = text
                elif i == 5: row_data['health'] = text
                elif i == 6: row_data['balance'] = text
                elif i == 7: row_data['provider'] = text
            accounts_info.append(row_data)
            print(f"  行数据: {row_data}")

        test_results["accounts"] = accounts_info

        # ========== 阶段4: 找到并操作云账号 ==========
        print("\n>>> 阶段4: 操作云账号")

        if len(table_rows) > 0:
            first_row = table_rows[0]

            # 查找行内的同步按钮
            sync_btn = first_row.locator('button:has-text("同步")').first

            if sync_btn:
                print("✓ 找到同步按钮")

                # 截图：同步前
                page.screenshot(path=os.path.join(OUTPUT_DIR, "before_sync.png"))

                # 点击同步按钮
                sync_btn.click()
                print("✓ 点击同步按钮")
                page.wait_for_timeout(2000)

                # 截图：同步对话框
                page.screenshot(path=os.path.join(OUTPUT_DIR, "sync_dialog.png"))

                # 检查是否出现了同步对话框
                sync_dialog = page.locator('.el-dialog:has-text("同步云账号")').first

                if sync_dialog:
                    print("✓ 同步对话框已打开")

                    # 选择同步模式（全量同步）
                    full_sync_radio = sync_dialog.locator('.el-radio:has-text("全量同步")').first
                    if full_sync_radio:
                        full_sync_radio.click()
                        print("✓ 选择全量同步")

                    # 选择全部资源类型
                    all_resources_radio = sync_dialog.locator('.el-radio:has-text("全部资源类型")').first
                    if all_resources_radio:
                        all_resources_radio.click()
                        print("✓ 选择全部资源类型")

                    page.wait_for_timeout(1000)

                    # 点击确认同步按钮
                    confirm_btn = sync_dialog.locator('button:has-text("确认同步")').first
                    if confirm_btn:
                        confirm_btn.click()
                        print("✓ 点击确认同步")

                        page.wait_for_timeout(5000)

                        # 截图：同步后
                        page.screenshot(path=os.path.join(OUTPUT_DIR, "after_sync.png"))

                        # 检查消息提示
                        messages = page.locator('.el-message').all()
                        for msg in messages:
                            text = msg.inner_text()
                            print(f"  消息: {text}")

                        test_results["sync_status"] = "success"
                    else:
                        print("✗ 未找到确认同步按钮")
                        test_results["sync_status"] = "no_confirm_button"
                else:
                    print("✗ 同步对话框未打开")
                    test_results["sync_status"] = "no_dialog"
            else:
                print("✗ 未找到同步按钮")
                test_results["sync_status"] = "no_sync_button"
        else:
            print("✗ 没有云账号")
            test_results["sync_status"] = "no_accounts"

        # ========== 阶段5: 检查定时任务 ==========
        print("\n>>> 阶段5: 检查定时任务")
        page.goto('http://localhost:3000/scheduled-tasks', wait_until='networkidle')
        page.wait_for_timeout(3000)
        page.screenshot(path=os.path.join(OUTPUT_DIR, "scheduled_tasks.png"))

        task_rows = page.locator('.el-table__row').all()
        tasks = []
        for row in task_rows:
            cells = row.locator('.el-table__cell').all()
            if cells:
                tasks.append(cells[0].inner_text())

        print(f"定时任务: {tasks}")
        test_results["scheduled_tasks"] = tasks

        # ========== 阶段6: 检查同步日志 ==========
        print("\n>>> 阶段6: 检查同步日志")
        page.goto('http://localhost:3000/cloud-management/sync-policies', wait_until='networkidle')
        page.wait_for_timeout(3000)
        page.screenshot(path=os.path.join(OUTPUT_DIR, "sync_policies.png"))

        # ========== 阶段7: 分析 API ==========
        print("\n>>> 阶段7: API 调用分析")
        print(f"云账号相关 API ({len(api_calls)} 个):")
        for api in api_calls:
            print(f"  {api['method']} {api['url']}")

        # ========== 阶段8: 控制台错误 ==========
        print("\n>>> 阶段8: 控制台错误")
        if console_errors:
            print(f"发现 {len(console_errors)} 个错误:")
            for err in console_errors[:5]:
                print(f"  {err[:100]}")
        else:
            print("✓ 无错误")

        # ========== 总结 ==========
        print("\n" + "="*70)
        print("测试完成")
        print("="*70)

        save_json(test_results, "sync_test_v2.json")
        save_json(api_calls, "sync_api_calls.json")

        page.wait_for_timeout(3000)
        browser.close()

if __name__ == '__main__':
    main()