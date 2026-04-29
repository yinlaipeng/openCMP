#!/usr/bin/env python3
"""
openCMP 多云管理模块全量测试脚本
测试页面功能、API响应、UI状态
"""

from playwright.sync_api import sync_playwright
import json
import os
import time

# 配置
BASE_URL = "http://localhost:3000"
API_URL = "http://localhost:8080"
USERNAME = "admin"
PASSWORD = "admin123"
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/test_results"

# 测试页面列表
PAGES_TO_TEST = [
    {
        "name": "云账户管理",
        "path": "/cloud-accounts",
        "expected_elements": ["表格", "工具栏", "搜索栏"],
        "api_endpoint": "/api/v1/cloud-accounts"
    },
    {
        "name": "同步策略",
        "path": "/cloud-management/sync-policies",
        "expected_elements": ["表格", "工具栏"],
        "api_endpoint": "/api/v1/sync-policies"
    },
    {
        "name": "云用户组",
        "path": "/cloud-management/cloud-user-groups",
        "expected_elements": ["表格", "工具栏", "筛选区"],
        "api_endpoint": "/api/v1/cloud-user-groups"
    },
    {
        "name": "代理管理",
        "path": "/cloud-management/proxies",
        "expected_elements": ["表格", "工具栏", "筛选区"],
        "api_endpoint": "/api/v1/proxies"
    },
    {
        "name": "定时同步任务",
        "path": "/scheduled-tasks",
        "expected_elements": ["表格", "对话框"],
        "api_endpoint": "/api/v1/scheduled-tasks"
    },
    {
        "name": "同步日志",
        "path": "/sync-logs",
        "expected_elements": ["统计卡片", "表格", "筛选区"],
        "api_endpoint": "/api/v1/sync-logs"
    }
]

def ensure_output_dir():
    """确保输出目录存在"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def login(page):
    """登录系统"""
    print("正在登录系统...")
    page.goto(f"{BASE_URL}/login", wait_until="networkidle")
    time.sleep(1)

    # 填写登录表单
    username_input = page.locator('input[type="text"]').first
    password_input = page.locator('input[type="password"]').first

    username_input.fill(USERNAME)
    password_input.fill(PASSWORD)

    # 点击登录按钮
    login_button = page.locator('button:has-text("登录")').first
    login_button.click()

    # 等待登录完成
    page.wait_for_load_state("networkidle")
    time.sleep(3)

    # 检查是否成功登录
    current_url = page.url
    if "/login" in current_url:
        print("警告: 登录可能失败，仍在登录页面")
        return False
    print(f"登录成功，当前URL: {current_url}")
    return True

def test_page(page, page_info):
    """测试单个页面"""
    print(f"\n--- 测试页面: {page_info['name']} ---")

    # 导航到页面
    page.goto(f"{BASE_URL}{page_info['path']}", wait_until="networkidle")
    time.sleep(2)

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{page_info['name'].replace('/', '_')}.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"截图保存: {screenshot_path}")

    # 检查页面元素
    results = {
        "name": page_info["name"],
        "path": page_info["path"],
        "url": page.url,
        "screenshot": screenshot_path,
        "checks": {}
    }

    # 1. 检查是否存在 EmptyState 或 "暂无" 文本
    empty_state_count = page.locator('text="暂无"').count()
    results["checks"]["empty_state_text"] = empty_state_count
    print(f"  '暂无'文本数量: {empty_state_count}")

    # 2. 检查表格是否存在
    table_count = page.locator('.el-table').count()
    results["checks"]["table_exists"] = table_count > 0
    print(f"  表格存在: {table_count > 0}")

    # 3. 检查表格表头
    if table_count > 0:
        headers = page.locator('.el-table__header th').all()
        header_texts = []
        for h in headers:
            cell = h.locator('.cell')
            if cell.count() > 0:
                txt = cell.inner_text().strip()
                if txt:
                    header_texts.append(txt)
        results["checks"]["table_headers"] = header_texts
        print(f"  表头数量: {len(header_texts)}")
        print(f"  表头内容: {header_texts[:5]}...")

    # 4. 检查工具栏按钮
    toolbar_buttons = page.locator('.toolbar button, .page-header button, .tool-bar button').all()
    button_texts = [b.inner_text().strip() for b in toolbar_buttons if b.inner_text().strip()]
    results["checks"]["toolbar_buttons"] = button_texts
    print(f"  工具栏按钮: {button_texts}")

    # 5. 检查筛选区
    filter_count = page.locator('.filter-card, .search-bar').count()
    results["checks"]["filter_area"] = filter_count > 0
    print(f"  筛选区存在: {filter_count > 0}")

    # 6. 检查分页组件
    pagination_count = page.locator('.el-pagination').count()
    results["checks"]["pagination"] = pagination_count > 0
    print(f"  分页组件存在: {pagination_count > 0}")

    # 7. 检查是否有错误提示
    error_count = page.locator('.el-message--error, text="服务器内部错误"').count()
    results["checks"]["has_error"] = error_count > 0
    print(f"  错误提示: {error_count > 0}")

    # 8. 统计卡片检查（针对同步日志）
    stats_count = page.locator('.stat-card, .stats-row').count()
    results["checks"]["stats_cards"] = stats_count
    print(f"  统计卡片数量: {stats_count}")

    return results

def main():
    """主测试流程"""
    ensure_output_dir()

    all_results = {
        "test_time": time.strftime("%Y-%m-%d %H:%M:%S"),
        "pages": []
    }

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(viewport={"width": 1920, "height": 1080})

        # 登录
        login_success = login(page)
        all_results["login_success"] = login_success

        if login_success:
            # 测试每个页面
            for page_info in PAGES_TO_TEST:
                try:
                    result = test_page(page, page_info)
                    all_results["pages"].append(result)
                except Exception as e:
                    print(f"测试页面 {page_info['name']} 失败: {e}")
                    all_results["pages"].append({
                        "name": page_info["name"],
                        "error": str(e)
                    })

        browser.close()

    # 保存测试结果
    results_path = os.path.join(OUTPUT_DIR, "test_results.json")
    with open(results_path, 'w', encoding='utf-8') as f:
        json.dump(all_results, f, ensure_ascii=False, indent=2)
    print(f"\n测试结果保存: {results_path}")

    # 生成测试报告
    print("\n" + "="*60)
    print("测试报告摘要")
    print("="*60)

    for page_result in all_results.get("pages", []):
        name = page_result.get("name", "未知")
        checks = page_result.get("checks", {})

        status = "✅ 正常" if checks.get("table_exists") and not checks.get("has_error") else "❌ 异常"
        print(f"\n{name}: {status}")

        if "error" in page_result:
            print(f"  错误: {page_result['error']}")
        else:
            print(f"  表格存在: {checks.get('table_exists', False)}")
            print(f"  表头数量: {len(checks.get('table_headers', []))}")
            print(f"  '暂无'文本: {checks.get('empty_state_text', 0)}")
            print(f"  工具栏按钮: {checks.get('toolbar_buttons', [])}")
            print(f"  有错误提示: {checks.get('has_error', False)}")

    print("\n" + "="*60)
    print("测试完成!")
    print("="*60)

if __name__ == "__main__":
    main()