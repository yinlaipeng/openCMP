#!/usr/bin/env python3
"""
openCMP 云账号添加和资源同步全量测试
测试流程:
1. 登录系统
2. 进入云账号页面
3. 检查云账号 "aliyun-test"
4. 测试同步全部资源功能
5. 检查 API 调用和响应
6. 验证定时任务是否注册
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
    test_results = {
        "login": {},
        "cloud_accounts_page": {},
        "account_list": {},
        "sync_test": {},
        "api_analysis": {},
        "scheduled_tasks": {},
        "console_errors": [],
        "summary": {}
    }

    api_calls = []
    responses = []
    console_errors = []

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        context = browser.new_context(viewport={'width': 1280, 'height': 900})
        page = context.new_page()

        # 监听 API 请求
        page.on('request', lambda req: api_calls.append({
            'url': req.url,
            'method': req.method,
            'type': req.resource_type,
            'body': req.post_data[:500] if req.post_data else None
        }) if '/api/' in req.url else None)

        # 监听响应
        page.on('response', lambda res: responses.append({
            'url': res.url,
            'status': res.status,
            'body': res.text()[:1000] if res.status < 400 else None
        }) if '/api/' in res.url else None)

        # 监听控制台错误
        page.on('console', lambda msg: console_errors.append({
            'type': msg.type,
            'text': msg.text
        }) if msg.type == 'error' else None)

        print("\n" + "="*70)
        print("openCMP 云账号添加和资源同步全量测试")
        print("="*70)

        # ========== 阶段1: 登录 ==========
        print("\n>>> 阶段1: 登录系统")
        try:
            page.goto('http://localhost:3000/login', wait_until='networkidle')
            page.wait_for_timeout(2000)

            page.locator('input').nth(0).fill('admin')
            page.locator('input').nth(1).fill('admin@123')
            page.locator('button.el-button--primary').first.click()

            page.wait_for_timeout(5000)

            if 'dashboard' in page.url or 'cloud-accounts' in page.url:
                print("✓ 登录成功")
                test_results["login"] = {"status": "success", "url": page.url}
            else:
                print(f"✗ 登录失败: {page.url}")
                test_results["login"] = {"status": "error", "url": page.url}
        except Exception as e:
            test_results["login"] = {"status": "error", "message": str(e)}
            print(f"✗ 登录异常: {e}")

        # ========== 阶段2: 进入云账号页面 ==========
        print("\n>>> 阶段2: 进入云账号页面")
        try:
            page.goto('http://localhost:3000/cloud-accounts', wait_until='networkidle')
            page.wait_for_timeout(3000)

            page.screenshot(path=os.path.join(OUTPUT_DIR, "cloud_accounts_page.png"))
            print("✓ 云账号页面加载成功")

            # 检查页面元素
            buttons = page.locator('button').all()
            table_rows = page.locator('.el-table__row').all()

            test_results["cloud_accounts_page"] = {
                "status": "success",
                "button_count": len(buttons),
                "table_row_count": len(table_rows),
                "title": page.title()
            }
            print(f"  按钮: {len(buttons)} 个, 表格行: {len(table_rows)} 行")
        except Exception as e:
            test_results["cloud_accounts_page"] = {"status": "error", "message": str(e)}
            print(f"✗ 页面加载异常: {e}")

        # ========== 阶段3: 检查云账号列表 ==========
        print("\n>>> 阶段3: 检查云账号列表")
        try:
            page.wait_for_timeout(2000)

            # 检查表格内容
            table_rows = page.locator('.el-table__row').all()
            account_names = []
            account_ids = []

            for row in table_rows:
                cells = row.locator('.el-table__cell').all()
                if len(cells) > 0:
                    name = cells[0].inner_text() if cells[0] else ''
                    account_names.append(name)

            # 检查是否存在 aliyun-test
            aliyun_test_exists = 'aliyun-test' in account_names
            print(f"  云账号列表: {account_names}")
            print(f"  aliyun-test 存在: {aliyun_test_exists}")

            test_results["account_list"] = {
                "status": "success",
                "accounts": account_names,
                "aliyun_test_exists": aliyun_test_exists,
                "total_count": len(account_names)
            }

            page.screenshot(path=os.path.join(OUTPUT_DIR, "account_list.png"))

            if not aliyun_test_exists:
                print("\n⚠️  云账号 'aliyun-test' 不存在，需要先添加")
                # TODO: 添加云账号流程

        except Exception as e:
            test_results["account_list"] = {"status": "error", "message": str(e)}
            print(f"✗ 检查异常: {e}")

        # ========== 阶段4: 查找并点击同步按钮 ==========
        print("\n>>> 阶段4: 测试同步全部资源功能")

        # 查找同步按钮
        sync_buttons = page.locator('button:has-text("同步")').all()
        sync_all_buttons = page.locator('button:has-text("同步全部资源")').all()

        print(f"  找到同步按钮: {len(sync_buttons)} 个")
        print(f"  找到同步全部资源按钮: {len(sync_all_buttons)} 个")

        if len(sync_all_buttons) > 0:
            print("\n>>> 点击 '同步全部资源' 按钮")
            try:
                sync_all_buttons[0].click()
                print("✓ 点击成功")

                # 等待响应
                page.wait_for_timeout(5000)

                # 检查是否有响应提示
                el_messages = page.locator('.el-message').all()
                print(f"  Element Plus 消息: {len(el_messages)} 个")
                for msg in el_messages:
                    text = msg.inner_text()
                    print(f"    - {text}")

                # 截图
                page.screenshot(path=os.path.join(OUTPUT_DIR, "sync_clicked.png"))

                test_results["sync_test"] = {
                    "status": "clicked",
                    "messages": [m.inner_text() for m in el_messages]
                }

            except Exception as e:
                test_results["sync_test"] = {"status": "error", "message": str(e)}
                print(f"✗ 点击异常: {e}")
        else:
            print("⚠️  未找到 '同步全部资源' 按钮")
            test_results["sync_test"] = {"status": "no_button_found"}

            # 打印所有按钮文本
            all_buttons = page.locator('button').all()
            print("\n  所有按钮文本:")
            for btn in all_buttons:
                text = btn.inner_text()
                print(f"    - '{text}'")

        # ========== 阶段5: 分析 API 调用 ==========
        print("\n>>> 阶段5: 分析 API 调用")

        # 筛选相关 API
        sync_apis = []
        account_apis = []
        task_apis = []

        for call in api_calls:
            url = call['url']
            if 'sync' in url.lower():
                sync_apis.append(call)
            if 'cloud-account' in url.lower() or 'account' in url.lower():
                account_apis.append(call)
            if 'task' in url.lower() or 'schedule' in url.lower():
                task_apis.append(call)

        print(f"\n  云账号相关 API ({len(account_apis)} 个):")
        for api in account_apis:
            print(f"    {api['method']} {api['url']}")

        print(f"\n  同步相关 API ({len(sync_apis)} 个):")
        for api in sync_apis:
            print(f"    {api['method']} {api['url']}")
            if api['body']:
                print(f"    Body: {api['body'][:200]}")

        print(f"\n  任务相关 API ({len(task_apis)} 个):")
        for api in task_apis:
            print(f"    {api['method']} {api['url']}")

        test_results["api_analysis"] = {
            "total_calls": len(api_calls),
            "sync_apis": sync_apis,
            "account_apis": [{"url": a['url'], "method": a['method']} for a in account_apis],
            "task_apis": [{"url": t['url'], "method": t['method']} for t in task_apis]
        }

        # ========== 阶段6: 检查定时任务页面 ==========
        print("\n>>> 阶段6: 检查定时任务是否注册")
        try:
            page.goto('http://localhost:3000/scheduled-tasks', wait_until='networkidle')
            page.wait_for_timeout(3000)

            page.screenshot(path=os.path.join(OUTPUT_DIR, "scheduled_tasks.png"))

            # 检查表格
            task_rows = page.locator('.el-table__row').all()
            task_names = []

            for row in task_rows:
                cells = row.locator('.el-table__cell').all()
                if len(cells) > 0:
                    name = cells[0].inner_text() if cells[0] else ''
                    task_names.append(name)

            print(f"  定时任务列表: {task_names}")

            test_results["scheduled_tasks"] = {
                "status": "success",
                "tasks": task_names,
                "total_count": len(task_names)
            }

        except Exception as e:
            test_results["scheduled_tasks"] = {"status": "error", "message": str(e)}
            print(f"✗ 定时任务页面异常: {e}")

        # ========== 阶段7: 检查控制台错误 ==========
        print("\n>>> 阶段7: 检查控制台错误")
        if console_errors:
            print(f"  发现 {len(console_errors)} 个控制台错误:")
            for err in console_errors:
                print(f"    [{err['type']}] {err['text'][:150]}")
        else:
            print("  ✓ 无控制台错误")

        test_results["console_errors"] = console_errors

        # ========== 阶段8: 直接测试后端 API ==========
        print("\n>>> 阶段8: 直接测试后端 API")

        # 获取 token
        token = page.evaluate("localStorage.getItem('token')")

        if token:
            print(f"  Token: {token[:50]}...")

            # 测试同步 API
            import requests as req_lib

            headers = {'Authorization': f'Bearer {token}'}

            # 获取云账号列表
            print("\n  测试 GET /api/v1/cloud-accounts")
            try:
                r = req_lib.get('http://localhost:8080/api/v1/cloud-accounts', headers=headers)
                print(f"    Status: {r.status_code}")
                if r.status_code == 200:
                    data = r.json()
                    accounts = data.get('items', []) or data.get('data', [])
                    print(f"    云账号数量: {len(accounts)}")
                    for acc in accounts[:3]:
                        print(f"      - {acc.get('name', 'N/A')} ({acc.get('provider', 'N/A')})")
            except Exception as e:
                print(f"    错误: {e}")

            # 测试同步 API
            print("\n  测试 POST /api/v1/cloud-accounts/sync-all")
            try:
                r = req_lib.post('http://localhost:8080/api/v1/cloud-accounts/sync-all', headers=headers)
                print(f"    Status: {r.status_code}")
                print(f"    Response: {r.text[:200] if r.text else 'N/A'}")
            except Exception as e:
                print(f"    错误: {e}")

            # 测试定时任务 API
            print("\n  测试 GET /api/v1/scheduled-tasks")
            try:
                r = req_lib.get('http://localhost:8080/api/v1/scheduled-tasks', headers=headers)
                print(f"    Status: {r.status_code}")
                if r.status_code == 200:
                    data = r.json()
                    tasks = data.get('items', []) or data.get('data', [])
                    print(f"    定时任务数量: {len(tasks)}")
            except Exception as e:
                print(f"    错误: {e}")

        # ========== 总结 ==========
        print("\n" + "="*70)
        print("测试总结")
        print("="*70)

        success_count = sum(1 for k, v in test_results.items()
                          if isinstance(v, dict) and v.get('status') in ['success', 'clicked'])

        test_results["summary"] = {
            "success_count": success_count,
            "total_phases": 7,
            "api_call_count": len(api_calls),
            "console_error_count": len(console_errors)
        }

        save_json(test_results, "test_results.json")
        save_json(api_calls, "api_calls.json")
        save_json(responses, "api_responses.json")

        print(f"\n成功阶段: {success_count}/7")
        print(f"API 调用: {len(api_calls)} 次")
        print(f"控制台错误: {len(console_errors)} 个")
        print(f"\n输出目录: {OUTPUT_DIR}")

        page.wait_for_timeout(3000)
        browser.close()

if __name__ == '__main__':
    main()