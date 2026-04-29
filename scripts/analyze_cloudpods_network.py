#!/usr/bin/env python3
"""
Cloudpods 网络服务页面分析脚本
分析页面：EIP、NAT网关、DNS解析、IPv6网关
"""

import json
import os
from playwright.sync_api import sync_playwright

# 配置
BASE_URL = "https://127.0.0.1"
LOGIN_URL = "https://127.0.0.1/login"
USERNAME = "admin"
PASSWORD = "admin@123"

# 要分析的页面
PAGES_TO_ANALYZE = [
    {"name": "EIP弹性公网IP", "url": "/eip", "module": "网络-网络服务"},
    {"name": "NAT网关", "url": "/nat", "module": "网络-网络服务"},
    {"name": "DNS解析", "url": "/vpc-network", "module": "网络-网络服务"},
    {"name": "IPv6网关", "url": "/ipv6-gateway", "module": "网络-网络服务"},
]

OUTPUT_DIR = "/tmp/cloudpods_analysis"

def setup_output_dir():
    """创建输出目录"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    for page in PAGES_TO_ANALYZE:
        page_dir = os.path.join(OUTPUT_DIR, page["name"])
        os.makedirs(page_dir, exist_ok=True)

def save_analysis(page_name, data):
    """保存分析结果"""
    filepath = os.path.join(OUTPUT_DIR, page_name, "analysis.json")
    with open(filepath, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=2)
    print(f"已保存分析结果: {filepath}")

def login(page):
    """登录到 Cloudpods"""
    print("正在登录...")
    page.goto(LOGIN_URL, wait_until="networkidle")
    page.wait_for_timeout(2000)

    # 检查是否已经在登录页面
    try:
        # 填写登录表单
        username_input = page.locator('input[type="text"], input[name="username"], input[placeholder*="用户"]').first
        password_input = page.locator('input[type="password"]').first

        if username_input.is_visible():
            username_input.fill(USERNAME)
            password_input.fill(PASSWORD)

            # 点击登录按钮
            login_btn = page.locator('button:has-text("登录"), button:has-text("Login"), button[type="submit"]').first
            login_btn.click()

            page.wait_for_load_state("networkidle")
            page.wait_for_timeout(3000)
            print("登录成功")
    except Exception as e:
        print(f"登录过程: {e}")

def analyze_page(page, page_info):
    """分析单个页面"""
    print(f"\n{'='*60}")
    print(f"分析页面: {page_info['name']} ({page_info['url']})")
    print(f"{'='*60}")

    url = BASE_URL + page_info["url"]
    page.goto(url, wait_until="networkidle")
    page.wait_for_timeout(3000)

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, page_info["name"], "screenshot.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"截图已保存: {screenshot_path}")

    # 保存 HTML
    html_path = os.path.join(OUTPUT_DIR, page_info["name"], "page.html")
    with open(html_path, 'w', encoding='utf-8') as f:
        f.write(page.content())
    print(f"HTML已保存: {html_path}")

    analysis = {
        "page_name": page_info["name"],
        "page_url": page_info["url"],
        "module": page_info["module"],
        "toolbar_buttons": [],
        "search_area": {},
        "table_columns": [],
        "operations": [],
        "api_calls": [],
        "create_dialog": {},
    }

    # 1. 分析工具栏按钮
    try:
        toolbar = page.locator('.page-toolbar, .toolbar, [class*="toolbar"]').first
        if toolbar.is_visible():
            buttons = toolbar.locator('button, .ant-btn, [role="button"]').all()
            for btn in buttons:
                try:
                    text = btn.inner_text() if btn.inner_text() else btn.get_attribute('aria-label') or ""
                    btn_type = btn.get_attribute('class') or ""
                    analysis["toolbar_buttons"].append({
                        "text": text.strip(),
                        "type": btn_type,
                        "is_primary": "primary" in btn_type,
                        "is_disabled": btn.is_disabled()
                    })
                except:
                    pass
    except Exception as e:
        print(f"工具栏分析: {e}")

    # 2. 分析搜索区域
    try:
        search_area = page.locator('.search-box-wrap, .search-bar, [class*="search"]').first
        if search_area.is_visible():
            inputs = search_area.locator('input, .ant-input').all()
            selects = search_area.locator('select, .ant-select').all()
            analysis["search_area"] = {
                "inputs": [{"placeholder": inp.get_attribute('placeholder') or ""} for inp in inputs],
                "selects_count": len(selects),
                "html": search_area.inner_html()[:500]
            }
    except Exception as e:
        print(f"搜索区分析: {e}")

    # 3. 分析表格列
    try:
        table = page.locator('.vxe-table, .ant-table, table').first
        if table.is_visible():
            headers = table.locator('th, .vxe-table--header .vxe-header--column').all()
            for header in headers:
                try:
                    text = header.inner_text().strip()
                    if text and text not in ['', ' ', '\n']:
                        analysis["table_columns"].append(text)
                except:
                    pass
    except Exception as e:
        print(f"表格分析: {e}")

    # 4. 分析操作列
    try:
        operation_col = page.locator('.vxe-table--body tr').first
        if operation_col.is_visible():
            ops = operation_col.locator('button, .ant-btn, a, [role="button"]').all()
            for op in ops:
                try:
                    text = op.inner_text().strip()
                    if text:
                        analysis["operations"].append(text)
                except:
                    pass
    except Exception as e:
        print(f"操作列分析: {e}")

    # 5. 尝试打开新建弹窗
    try:
        create_btn = page.locator('button:has-text("新建"), button:has-text("创建"), button:has-text("申请"), .ant-btn-primary').first
        if create_btn.is_visible() and not create_btn.is_disabled():
            create_btn.click()
            page.wait_for_timeout(2000)

            # 截图弹窗
            dialog_screenshot = os.path.join(OUTPUT_DIR, page_info["name"], "create_dialog.png")
            page.screenshot(path=dialog_screenshot)

            # 分析弹窗字段
            dialog = page.locator('.ant-modal, .el-dialog, [role="dialog"]').first
            if dialog.is_visible():
                form_items = dialog.locator('.ant-form-item, .el-form-item').all()
                dialog_fields = []
                for item in form_items:
                    try:
                        label = item.locator('label, .ant-form-item-label').inner_text().strip()
                        input_type = ""
                        inp = item.locator('input').first
                        if inp.is_visible():
                            input_type = "input"
                        sel = item.locator('select, .ant-select').first
                        if sel.is_visible():
                            input_type = "select"
                        textarea = item.locator('textarea').first
                        if textarea.is_visible():
                            input_type = "textarea"
                        dialog_fields.append({
                            "label": label,
                            "input_type": input_type
                        })
                    except:
                        pass

                analysis["create_dialog"] = {
                    "fields": dialog_fields,
                    "has_tabs": dialog.locator('.ant-tabs, .el-tabs').count() > 0
                }

                # 关闭弹窗
                close_btn = dialog.locator('button:has-text("取消"), button:has-text("Cancel"), .ant-modal-close').first
                if close_btn.is_visible():
                    close_btn.click()
                    page.wait_for_timeout(1000)
    except Exception as e:
        print(f"新建弹窗分析: {e}")

    save_analysis(page_info["name"], analysis)
    return analysis

def main():
    setup_output_dir()

    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--ignore-certificate-errors-spki-list']
        )
        context = browser.new_context(
            ignore_https_errors=True,
            record_har_path=os.path.join(OUTPUT_DIR, "network_requests.har")
        )
        page = context.new_page()

        # 登录
        login(page)

        # 分析每个页面
        all_results = []
        for page_info in PAGES_TO_ANALYZE:
            result = analyze_page(page, page_info)
            all_results.append(result)

        browser.close()

        # 输出汇总
        print("\n" + "="*60)
        print("分析汇总")
        print("="*60)
        for result in all_results:
            print(f"\n页面: {result['page_name']}")
            print(f"  工具栏按钮: {len(result['toolbar_buttons'])} 个")
            print(f"  表格列: {result['table_columns']}")
            print(f"  操作按钮: {result['operations']}")
            if result['create_dialog'].get('fields'):
                print(f"  新建弹窗字段: {result['create_dialog']['fields']}")

if __name__ == "__main__":
    main()