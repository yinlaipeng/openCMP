#!/usr/bin/env python3
"""
Cloudpods 新建弹窗详细分析
"""

import json
import os
from playwright.sync_api import sync_playwright

BASE_URL = "https://127.0.0.1"
USERNAME = "admin"
PASSWORD = "admin@123"

PAGES = [
    {"name": "EIP", "url": "/eip", "create_btn": "Create"},
    {"name": "NAT", "url": "/nat", "create_btn": "Create"},
]

OUTPUT_DIR = "/tmp/cloudpods_net"

def login(page):
    page.goto(BASE_URL + "/login", wait_until="networkidle", timeout=60000)
    page.wait_for_timeout(3000)
    page.locator('input[placeholder="Please enter your username"]').fill(USERNAME)
    page.locator('input[placeholder="Please enter your password"]').fill(PASSWORD)
    page.locator('button[type="submit"]').click()
    page.wait_for_load_state("networkidle", timeout=30000)
    page.wait_for_timeout(5000)

def analyze_create_dialog(page, page_info):
    print(f"\n{'='*50}")
    print(f"分析新建弹窗: {page_info['name']}")
    print(f"{'='*50}")

    page.goto(BASE_URL + page_info["url"], wait_until="networkidle", timeout=60000)
    page.wait_for_timeout(3000)

    # 点击新建按钮
    create_btn = page.locator(f'button:has-text("{page_info["create_btn"]}")')
    if create_btn.count() > 0:
        create_btn.first.click()
        page.wait_for_timeout(3000)

        # 截图弹窗
        screenshot_path = os.path.join(OUTPUT_DIR, page_info["name"], "create_dialog.png")
        page.screenshot(path=screenshot_path, full_page=True)
        print(f"弹窗截图: {screenshot_path}")

        # 分析弹窗
        dialog = page.locator('.ant-modal-content, [role="dialog"]').first
        if dialog.is_visible():
            result = {
                "title": "",
                "tabs": [],
                "form_sections": {},
                "all_fields": []
            }

            # 获取标题
            try:
                title = dialog.locator('.ant-modal-title').inner_text().strip()
                result["title"] = title
            except:
                pass

            # 获取 Tabs（如果有）
            try:
                tabs = dialog.locator('.ant-tabs-tab').all()
                for tab in tabs:
                    result["tabs"].append(tab.inner_text().strip())
            except:
                pass

            # 分析表单项
            try:
                form_items = dialog.locator('.ant-form-item').all()
                fields = []
                for item in form_items:
                    try:
                        label = item.locator('.ant-form-item-label label').inner_text().strip()
                        if not label:
                            continue

                        # 判断输入类型
                        field_type = "unknown"

                        # 检查是否有 select
                        if item.locator('.ant-select').count() > 0:
                            field_type = "select"
                            # 尝试获取选项
                            try:
                                item.locator('.ant-select').click()
                                page.wait_for_timeout(500)
                                options = page.locator('.ant-select-dropdown .ant-select-item').all()
                                opts_text = [opt.inner_text().strip() for opt in options[:10]]
                                fields.append({
                                    "label": label,
                                    "type": "select",
                                    "options": opts_text
                                })
                                # 关闭下拉
                                page.keyboard.press("Escape")
                                page.wait_for_timeout(500)
                            except:
                                fields.append({"label": label, "type": "select"})

                        elif item.locator('input[type="number"]').count() > 0:
                            field_type = "number"
                            placeholder = item.locator('input').get_attribute("placeholder") or ""
                            fields.append({"label": label, "type": "number", "placeholder": placeholder})

                        elif item.locator('input[type="text"]').count() > 0:
                            field_type = "text"
                            placeholder = item.locator('input').get_attribute("placeholder") or ""
                            fields.append({"label": label, "type": "text", "placeholder": placeholder})

                        elif item.locator('textarea').count() > 0:
                            field_type = "textarea"
                            placeholder = item.locator('textarea').get_attribute("placeholder") or ""
                            fields.append({"label": label, "type": "textarea", "placeholder": placeholder})

                        elif item.locator('.ant-input-number').count() > 0:
                            field_type = "number"
                            fields.append({"label": label, "type": "number"})

                        elif item.locator('.ant-radio-group').count() > 0:
                            field_type = "radio"
                            radios = item.locator('.ant-radio-wrapper').all()
                            radio_options = [r.inner_text().strip() for r in radios]
                            fields.append({"label": label, "type": "radio", "options": radio_options})

                        elif item.locator('.ant-checkbox-group').count() > 0:
                            field_type = "checkbox"
                            checkboxes = item.locator('.ant-checkbox-wrapper').all()
                            checkbox_options = [c.inner_text().strip() for c in checkboxes]
                            fields.append({"label": label, "type": "checkbox", "options": checkbox_options})

                        elif field_type == "unknown":
                            fields.append({"label": label, "type": "unknown"})

                    except Exception as e:
                        print(f"字段分析错误: {e}")

                result["all_fields"] = fields
            except Exception as e:
                print(f"表单分析错误: {e}")

            # 保存结果
            result_path = os.path.join(OUTPUT_DIR, page_info["name"], "create_dialog.json")
            with open(result_path, 'w', encoding='utf-8') as f:
                json.dump(result, f, ensure_ascii=False, indent=2)
            print(f"结果: {result_path}")

            # 打印字段
            print("\n新建弹窗字段:")
            for field in result["all_fields"]:
                print(f"  {field['label']}: {field['type']}")
                if "options" in field:
                    print(f"    选项: {field['options']}")

            # 关闭弹窗
            try:
                cancel_btn = page.locator('.ant-modal-close, button:has-text("Cancel")')
                if cancel_btn.count() > 0:
                    cancel_btn.first.click()
                    page.wait_for_timeout(1000)
            except:
                pass

            return result

    return {}

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--ignore-certificate-errors-spki-list']
        )
        context = browser.new_context(ignore_https_errors=True)
        page = context.new_page()

        login(page)

        for page_info in PAGES:
            analyze_create_dialog(page, page_info)

        browser.close()

if __name__ == "__main__":
    main()