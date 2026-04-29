#!/usr/bin/env python3
"""
深入分析 Cloudpods WAF策略和应用程序服务页面 - 操作列和新建弹窗

分析：
1. 操作列下拉菜单内容
2. 新建按钮和弹窗设计
3. API 数据结构
"""

import json
import os
import time
import asyncio
from playwright.sync_api import sync_playwright

OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/.specify/cloudpods_analysis"

def login_cloudpods(page, base_url: str, username: str, password: str):
    """登录 Cloudpods 系统"""
    login_url = f"{base_url}/auth/login"
    print(f"[INFO] 正在访问登录页面: {login_url}")

    page.goto(login_url, wait_until="networkidle", timeout=60000)
    time.sleep(3)

    # 输入登录凭证
    try:
        username_input = page.locator('input.ant-input[type="text"]').first
        password_input = page.locator('input.ant-input[type="password"]').first

        if username_input.is_visible() and password_input.is_visible():
            username_input.fill(username)
            time.sleep(0.5)
            password_input.fill(password)
            time.sleep(0.5)

            login_button = page.locator('button[type="submit"]').first
            if login_button.is_visible():
                login_button.click()
                time.sleep(5)
                page.wait_for_load_state("networkidle", timeout=30000)

                if "/login" not in page.url:
                    print(f"[INFO] 登录成功!")
                    return True
    except Exception as e:
        print(f"[ERROR] 登录异常: {e}")

    return False

def analyze_operations_dropdown(page, page_url: str, page_name: str):
    """分析操作列下拉菜单"""
    print(f"\n{'='*60}")
    print(f"[INFO] 分析操作列下拉菜单: {page_name}")
    print(f"{'='*60}")

    result = {
        "page_name": page_name,
        "operations_dropdown": []
    }

    # 导航到页面
    page.goto(page_url, wait_until="networkidle", timeout=60000)
    time.sleep(3)

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{page_name}_operations.png")
    page.screenshot(path=screenshot_path)
    print(f"[INFO] 截图: {screenshot_path}")

    # 检查是否有数据行（如果没有，无法分析操作列）
    try:
        rows = page.locator('.vxe-body--row, .ant-table-tbody tr').all()
        print(f"[INFO] 表格行数: {len(rows)}")

        if len(rows) > 0:
            # 尝试点击第一行的操作下拉菜单
            operations_cell = rows[0].locator('td:last-child, .vxe-body--column:last-child').first

            # 查找下拉按钮
            dropdown_btn = operations_cell.locator('.ant-dropdown-trigger, .vxe-dropdown--trigger, button:has(.anticon-down)').first
            if dropdown_btn.is_visible():
                print(f"[INFO] 找到操作下拉按钮，点击...")
                dropdown_btn.click()
                time.sleep(1)

                # 获取下拉菜单内容
                dropdown_menu = page.locator('.ant-dropdown-menu, .vxe-dropdown--panel').first
                if dropdown_menu.is_visible():
                    menu_items = dropdown_menu.locator('.ant-dropdown-menu-item, .vxe-dropdown--item').all()
                    for item in menu_items:
                        text = item.text_content().strip()
                        if text:
                            result["operations_dropdown"].append(text)
                            print(f"  - 操作: {text}")

                    # 关闭下拉菜单
                    page.keyboard.press("Escape")
                    time.sleep(0.5)
        else:
            print(f"[INFO] 无数据行，无法分析操作列下拉菜单")
            result["operations_dropdown"] = ["无法分析（无数据）"]
    except Exception as e:
        print(f"[WARN] 操作列分析异常: {e}")
        result["operations_dropdown"] = ["分析异常"]

    return result

def analyze_create_button(page, page_url: str, page_name: str):
    """分析新建按钮和弹窗"""
    print(f"\n{'='*60}")
    print(f"[INFO] 分析新建按钮: {page_name}")
    print(f"{'='*60}")

    result = {
        "page_name": page_name,
        "create_button_found": False,
        "create_modal": {}
    }

    # 导航到页面
    page.goto(page_url, wait_until="networkidle", timeout=60000)
    time.sleep(3)

    # 查找新建按钮
    try:
        # 查找页面头部区域的按钮
        header_area = page.locator('.page-header, .page-toolbar, header').first
        if header_area.is_visible():
            # 检查所有按钮，找到可能的"新建"按钮
            buttons = header_area.locator('button, .ant-btn').all()
            print(f"[INFO] 头部区域按钮数: {len(buttons)}")

            for btn in buttons:
                btn_text = btn.text_content().strip()
                btn_class = btn.get_attribute("class") or ""
                print(f"  按钮: '{btn_text}' class='{btn_class[:30]}'")

                # 查找带有特定class的按钮（通常是新建按钮）
                if "primary" in btn_class.lower() or "create" in btn_text.lower() or "新建" in btn_text:
                    print(f"[INFO] 可能的新建按钮: '{btn_text}'")
                    result["create_button_found"] = True
                    result["create_button_text"] = btn_text

                    # 如果是icon按钮（无文本），可能是新建按钮
                    if not btn_text:
                        print(f"[INFO] icon按钮（可能是新建）")

        # 尝试查找其他位置的新建按钮
        create_selectors = [
            'button.ant-btn-primary',
            '.ant-btn.ant-btn-primary',
            'button[type="button"].ant-btn-primary',
        ]

        for selector in create_selectors:
            btn = page.locator(selector).first
            if btn.is_visible():
                btn_text = btn.text_content().strip()
                print(f"[INFO] 找到 primary 按钮: selector={selector}, text='{btn_text}'")
                result["create_button_found"] = True
                if btn_text:
                    result["create_button_text"] = btn_text

                # 点击按钮
                btn.click()
                time.sleep(2)

                # 截取弹窗截图
                modal_screenshot = os.path.join(OUTPUT_DIR, f"{page_name}_create_modal.png")
                page.screenshot(path=modal_screenshot)
                print(f"[INFO] 弹窗截图: {modal_screenshot}")

                # 分析弹窗
                modal = page.locator('.ant-modal-content, .ant-modal').first
                if modal.is_visible():
                    # 弹窗标题
                    modal_title = modal.locator('.ant-modal-title, .ant-modal-header').first
                    if modal_title.is_visible():
                        title = modal_title.text_content().strip()
                        result["create_modal"]["title"] = title
                        print(f"[INFO] 弹窗标题: {title}")

                    # 弹窗表单字段
                    form = modal.locator('.ant-form').first
                    if form.is_visible():
                        form_items = form.locator('.ant-form-item').all()
                        fields = []
                        for item in form_items:
                            try:
                                label_elem = item.locator('.ant-form-item-label label').first
                                if label_elem.is_visible():
                                    label = label_elem.text_content().strip()
                                    required = item.locator('.ant-form-item-required').is_visible()

                                    # 检查输入类型
                                    input_elem = item.locator('input, .ant-select, .ant-radio-group').first
                                    input_type = "input"
                                    if input_elem.is_visible():
                                        input_class = input_elem.get_attribute("class") or ""
                                        if "ant-select" in input_class:
                                            input_type = "select"
                                        elif "ant-radio-group" in input_class:
                                            input_type = "radio"

                                    fields.append({
                                        "label": label,
                                        "required": required,
                                        "type": input_type
                                    })
                            except:
                                continue

                        result["create_modal"]["fields"] = fields
                        print(f"[INFO] 表单字段数: {len(fields)}")
                        for field in fields:
                            print(f"  - {field['label']} (required={field['required']}, type={field['type']})")

                    # 底部按钮
                    footer = modal.locator('.ant-modal-footer').first
                    if footer.is_visible():
                        footer_buttons = footer.locator('button').all()
                        footer_texts = [btn.text_content().strip() for btn in footer_buttons]
                        result["create_modal"]["footer_buttons"] = footer_texts
                        print(f"[INFO] 底部按钮: {footer_texts}")

                # 关闭弹窗
                close_btn = page.locator('.ant-modal-close').first
                if close_btn.is_visible():
                    close_btn.click()
                    time.sleep(1)

                break

    except Exception as e:
        print(f"[WARN] 新建按钮分析异常: {e}")

    return result

def capture_api_response(page, page_url: str, page_name: str):
    """捕获 API 响应"""
    print(f"\n{'='*60}")
    print(f"[INFO] 捕获 API 响应: {page_name}")
    print(f"{'='*60}")

    result = {
        "page_name": page_name,
        "api_responses": []
    }

    # 设置响应捕获
    responses_list = []

    def handle_response(response):
        url = response.url
        if "/api/" in url and ("waf" in url.lower() or "webapp" in url.lower()):
            try:
                json_body = response.json()
                responses_list.append({
                    "url": url,
                    "status": response.status,
                    "body": json_body
                })
            except:
                pass

    page.on("response", handle_response)

    # 导航到页面
    page.goto(page_url, wait_until="networkidle", timeout=60000)
    time.sleep(3)

    result["api_responses"] = responses_list

    # 打印捕获的响应
    for resp in responses_list:
        print(f"[INFO] API响应: {resp['url']}")
        print(f"  状态: {resp['status']}")
        if resp['body']:
            # 打印部分响应内容
            body_str = json.dumps(resp['body'], ensure_ascii=False)[:500]
            print(f"  响应内容(部分): {body_str}")

    return result

def main():
    """主函数"""
    base_url = "https://127.0.0.1"
    username = "admin"
    password = "admin@123"

    pages_to_analyze = [
        {"url": "/waf", "name": "waf_policy"},
        {"url": "/webapp", "name": "webapp_service"},
    ]

    os.makedirs(OUTPUT_DIR, exist_ok=True)

    print(f"[INFO] 深入分析 Cloudpods WAF/Webapp 页面...")

    all_results = {
        "operations": [],
        "create_buttons": [],
        "api_responses": []
    }

    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--disable-web-security']
        )

        context = browser.new_context(
            ignore_https_errors=True,
            viewport={"width": 1920, "height": 1080}
        )

        page = context.new_page()

        try:
            # 登录
            login_success = login_cloudpods(page, base_url, username, password)

            if login_success:
                for page_info in pages_to_analyze:
                    page_url = f"{base_url}{page_info['url']}"

                    # 分析操作列
                    ops_result = analyze_operations_dropdown(page, page_url, page_info['name'])
                    all_results["operations"].append(ops_result)

                    # 分析新建按钮
                    create_result = analyze_create_button(page, page_url, page_info['name'])
                    all_results["create_buttons"].append(create_result)

                    # 捕获 API 响应
                    api_result = capture_api_response(page, page_url, page_info['name'])
                    all_results["api_responses"].append(api_result)

            # 保存结果
            output_file = os.path.join(OUTPUT_DIR, "detailed_analysis.json")
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(all_results, f, ensure_ascii=False, indent=2)
            print(f"\n[INFO] 详细分析结果保存: {output_file}")

        except Exception as e:
            print(f"[ERROR] 分析异常: {e}")
            import traceback
            traceback.print_exc()

        finally:
            browser.close()

    print(f"[INFO] 分析完成!")

if __name__ == "__main__":
    main()