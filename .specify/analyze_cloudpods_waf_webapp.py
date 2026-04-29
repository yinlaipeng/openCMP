#!/usr/bin/env python3
"""
分析 Cloudpods WAF策略和应用程序服务页面

使用 Playwright 登录 Cloudpods 系统，分析：
1. WAF策略页面 (/waf)
2. 应用程序服务页面 (/webapp)

记录页面结构、工具栏、表格、弹窗设计和 API 接口。
"""

import json
import os
import sys
import time
from playwright.sync_api import sync_playwright

# 输出目录
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/.specify/cloudpods_analysis"

def setup_output_dir():
    """创建输出目录"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def login_cloudpods(page, base_url: str, username: str, password: str):
    """登录 Cloudpods 系统"""
    login_url = f"{base_url}/login"
    print(f"[INFO] 正在访问登录页面: {login_url}")

    page.goto(login_url, wait_until="networkidle", timeout=60000)
    time.sleep(2)

    # 检查是否已经在登录页面
    # 尝试查找登录表单
    try:
        # 查找用户名输入框 - Ant Design Vue 登录表单
        username_input = page.locator('input[type="text"], input[placeholder*="用户"], input[placeholder*="username"], #username, .ant-input[type="text"]').first
        password_input = page.locator('input[type="password"], input[placeholder*="密码"], input[placeholder*="password"], #password, .ant-input-password input').first

        if username_input.is_visible() and password_input.is_visible():
            print(f"[INFO] 找到登录表单，正在输入凭证...")
            username_input.fill(username)
            password_input.fill(password)
            time.sleep(1)

            # 点击登录按钮
            login_button = page.locator('button[type="submit"], button:has-text("登录"), button:has-text("Login"), .ant-btn-primary').first
            if login_button.is_visible():
                login_button.click()
                print(f"[INFO] 点击登录按钮...")
                time.sleep(3)
                page.wait_for_load_state("networkidle", timeout=30000)
                print(f"[INFO] 登录成功，当前URL: {page.url}")
                return True
    except Exception as e:
        print(f"[WARN] 登录流程异常: {e}")

    # 检查是否已经登录成功（通过检查页面元素）
    current_url = page.url
    if "/login" not in current_url:
        print(f"[INFO] 可能已登录成功，当前URL: {current_url}")
        return True

    return False

def analyze_page(page, page_url: str, page_name: str):
    """分析单个页面的结构"""
    print(f"\n{'='*60}")
    print(f"[INFO] 正在分析页面: {page_name} ({page_url})")
    print(f"{'='*60}")

    result = {
        "page_name": page_name,
        "page_url": page_url,
        "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
        "structure": {},
        "toolbar": [],
        "search": {},
        "table": {},
        "api_calls": [],
        "modals": {}
    }

    # 导航到页面
    page.goto(page_url, wait_until="networkidle", timeout=60000)
    time.sleep(2)

    # 截取全页截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{page_name}_full.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"[INFO] 截图保存: {screenshot_path}")

    # 1. 分析页面标题
    try:
        page_header = page.locator('.page-header h2, .page-header h3, h2.page-title, h3.page-title, .ant-page-header-title').first
        if page_header.is_visible():
            result["structure"]["page_title"] = page_header.text_content()
            print(f"[INFO] 页面标题: {result['structure']['page_title']}")
    except:
        result["structure"]["page_title"] = page_name

    # 2. 分析工具栏按钮
    toolbar_buttons = []
    try:
        # Ant Design Vue 工具栏按钮
        toolbar = page.locator('.page-toolbar, .toolbar, .ant-space').first
        if toolbar.is_visible():
            buttons = toolbar.locator('button, .ant-btn').all()
            for btn in buttons:
                btn_text = btn.text_content().strip()
                btn_type = btn.get_attribute("class") or ""
                is_disabled = btn.is_disabled()
                toolbar_buttons.append({
                    "text": btn_text,
                    "type": btn_type,
                    "disabled": is_disabled
                })
            print(f"[INFO] 工具栏按钮: {len(toolbar_buttons)} 个")
            for btn in toolbar_buttons:
                print(f"  - {btn['text']} ({btn['type']}) disabled={btn['disabled']}")
    except Exception as e:
        print(f"[WARN] 工具栏分析异常: {e}")
    result["toolbar"] = toolbar_buttons

    # 3. 分析搜索区域
    try:
        search_area = page.locator('.search-box-wrap, .search-bar, .filter-card, .ant-card').first
        if search_area.is_visible():
            # 搜索输入框
            search_input = search_area.locator('input[type="text"], .ant-input').first
            if search_input.is_visible():
                placeholder = search_input.get_attribute("placeholder") or ""
                result["search"]["input_placeholder"] = placeholder
                print(f"[INFO] 搜索框 placeholder: {placeholder}")

            # 筛选下拉框
            selects = search_area.locator('.ant-select, el-select').all()
            result["search"]["filters"] = len(selects)
            print(f"[INFO] 筛选下拉框: {len(selects)} 个")
    except Exception as e:
        print(f"[WARN] 搜索区域分析异常: {e}")

    # 4. 分析表格结构
    try:
        # VXE Table 或 Ant Design Table
        table = page.locator('.vxe-table, .ant-table, el-table').first
        if table.is_visible():
            # 表头列
            headers = table.locator('th, .vxe-header--column, .ant-table-thead th').all()
            header_texts = [h.text_content().strip() for h in headers]
            result["table"]["headers"] = header_texts
            print(f"[INFO] 表格表头: {header_texts}")

            # 表格行数
            rows = table.locator('tbody tr, .vxe-body--row, .ant-table-tbody tr').all()
            result["table"]["row_count"] = len(rows)
            print(f"[INFO] 表格行数: {len(rows)}")

            # 分页
            pagination = page.locator('.vxe-pager, .ant-pagination, el-pagination').first
            if pagination.is_visible():
                result["table"]["has_pagination"] = True
                print(f"[INFO] 有分页组件")
    except Exception as e:
        print(f"[WARN] 表格分析异常: {e}")

    # 5. 捕获 API 调用（通过检查页面内容推断）
    # 这里我们无法直接捕获网络请求，因为请求已经在 networkidle 之前完成
    # 但我们可以通过页面数据推断 API

    # 6. 尝试打开新建弹窗分析
    try:
        create_btn = page.locator('button:has-text("新建"), button:has-text("Create"), button:has-text("添加"), .ant-btn-primary:has-text("新建")').first
        if create_btn.is_visible() and not create_btn.is_disabled():
            print(f"[INFO] 找到新建按钮，尝试点击...")
            create_btn.click()
            time.sleep(2)

            # 截取弹窗截图
            modal_screenshot = os.path.join(OUTPUT_DIR, f"{page_name}_create_modal.png")
            page.screenshot(path=modal_screenshot)
            print(f"[INFO] 弹窗截图保存: {modal_screenshot}")

            # 分析弹窗结构
            modal = page.locator('.ant-modal, el-dialog, .vxe-modal').first
            if modal.is_visible():
                # 弹窗标题
                modal_title = modal.locator('.ant-modal-header, .ant-modal-title, el-dialog__header, .vxe-modal--header').first
                if modal_title.is_visible():
                    result["modals"]["create_title"] = modal_title.text_content().strip()
                    print(f"[INFO] 弹窗标题: {result['modals']['create_title']}")

                # 弹窗表单字段
                form_items = modal.locator('.ant-form-item, el-form-item').all()
                form_fields = []
                for item in form_items:
                    label = item.locator('.ant-form-item-label, el-form-item__label').first
                    if label.is_visible():
                        field_label = label.text_content().strip()
                        # 检查是否必填
                        required = item.locator('.ant-form-item-required, .required').is_visible()
                        form_fields.append({
                            "label": field_label,
                            "required": required
                        })
                result["modals"]["create_fields"] = form_fields
                print(f"[INFO] 弹窗表单字段: {len(form_fields)} 个")
                for field in form_fields:
                    print(f"  - {field['label']} (required={field['required']})")

                # 关闭弹窗
                close_btn = modal.locator('.ant-modal-close, el-dialog__close, .vxe-modal--close').first
                if close_btn.is_visible():
                    close_btn.click()
                    time.sleep(1)
    except Exception as e:
        print(f"[WARN] 弹窗分析异常: {e}")

    return result

def main():
    """主函数"""
    # 配置
    base_url = "https://127.0.0.1"
    username = "admin"
    password = "admin@123"

    # 要分析的页面
    pages_to_analyze = [
        {"url": "/waf", "name": "waf_policy"},
        {"url": "/webapp", "name": "webapp_service"},
    ]

    setup_output_dir()

    print(f"[INFO] 启动 Playwright 分析...")
    print(f"[INFO] 目标系统: {base_url}")
    print(f"[INFO] 输出目录: {OUTPUT_DIR}")

    results = []

    with sync_playwright() as p:
        # 启动浏览器（忽略 SSL 证书）
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--ignore-certificate-errors-spki-list']
        )

        # 创建浏览器上下文（忽略 SSL）
        context = browser.new_context(
            ignore_https_errors=True,
            viewport={"width": 1920, "height": 1080}
        )

        page = context.new_page()

        try:
            # 登录
            login_success = login_cloudpods(page, base_url, username, password)
            if not login_success:
                print(f"[ERROR] 登录失败，尝试继续分析...")

            # 分析各个页面
            for page_info in pages_to_analyze:
                page_url = f"{base_url}{page_info['url']}"
                result = analyze_page(page, page_url, page_info['name'])
                results.append(result)

            # 保存分析结果
            output_file = os.path.join(OUTPUT_DIR, "analysis_results.json")
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(results, f, ensure_ascii=False, indent=2)
            print(f"\n[INFO] 分析结果保存: {output_file}")

        except Exception as e:
            print(f"[ERROR] 分析过程异常: {e}")
            import traceback
            traceback.print_exc()

        finally:
            browser.close()

    print(f"\n[INFO] 分析完成!")
    print(f"[INFO] 截图保存在: {OUTPUT_DIR}")

if __name__ == "__main__":
    main()