#!/usr/bin/env python3
"""
Cloudpods 网络服务页面详细分析脚本 - 使用正确的选择器
"""

import json
import os
import time
from playwright.sync_api import sync_playwright

BASE_URL = "https://127.0.0.1"
USERNAME = "admin"
PASSWORD = "admin@123"

PAGES = [
    {"name": "EIP", "url": "/eip"},
    {"name": "NAT", "url": "/nat"},
    {"name": "DNS", "url": "/vpc-network"},
    {"name": "IPv6Gateway", "url": "/ipv6-gateway"},
]

OUTPUT_DIR = "/tmp/cloudpods_net"

def setup():
    for p in PAGES:
        os.makedirs(os.path.join(OUTPUT_DIR, p["name"]), exist_ok=True)

def login(page):
    print("导航到登录页面...")
    page.goto(BASE_URL + "/login", wait_until="networkidle", timeout=60000)
    page.wait_for_timeout(3000)

    # 使用正确的 placeholder
    username_input = page.locator('input[placeholder="Please enter your username"]')
    password_input = page.locator('input[placeholder="Please enter your password"]')

    print("填写登录信息...")
    username_input.fill(USERNAME)
    password_input.fill(PASSWORD)
    page.wait_for_timeout(500)

    # 点击登录按钮
    submit_btn = page.locator('button[type="submit"]')
    submit_btn.click()

    print("等待登录完成...")
    page.wait_for_load_state("networkidle", timeout=30000)
    page.wait_for_timeout(5000)

    # 检查登录是否成功
    current_url = page.url
    print(f"当前 URL: {current_url}")

    if "login" not in current_url:
        print("登录成功!")
        return True
    else:
        print("登录可能失败，尝试检查...")
        return False

def capture_api_calls(page):
    """捕获 API 调用"""
    api_calls = []

    def on_response(response):
        url = response.url
        if "/api/" in url and response.status == 200:
            try:
                body = response.text()
                api_calls.append({
                    "url": url,
                    "status": response.status,
                    "method": response.request.method,
                    "body_preview": body[:500] if body else None
                })
            except:
                api_calls.append({
                    "url": url,
                    "status": response.status,
                    "method": response.request.method,
                })

    page.on("response", on_response)
    return api_calls

def analyze_page(page, page_info, api_calls):
    print(f"\n{'='*50}")
    print(f"分析: {page_info['name']} ({page_info['url']})")
    print(f"{'='*50}")

    url = BASE_URL + page_info["url"]
    page.goto(url, wait_until="networkidle", timeout=60000)
    page.wait_for_timeout(3000)

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, page_info["name"], "page.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"截图: {screenshot_path}")

    # 保存 HTML
    html_path = os.path.join(OUTPUT_DIR, page_info["name"], "page.html")
    with open(html_path, 'w', encoding='utf-8') as f:
        f.write(page.content())
    print(f"HTML: {html_path}")

    analysis = {
        "page": page_info["name"],
        "url": page.url,
        "api_calls": [],
        "toolbar_buttons": [],
        "search_inputs": [],
        "table_columns": [],
        "tabs": [],
        "create_dialog": {},
    }

    # 关联 API 调用
    for call in api_calls:
        url_path = page_info["url"].replace("/", "")
        if url_path in call["url"] or call["url"].endswith(page_info["url"]):
            analysis["api_calls"].append(call)

    # 分析 Tabs
    try:
        tabs = page.locator('.ant-tabs-tab-pane-active, [role="tab"]').all()
        for tab in tabs[:10]:  # 限制数量
            text = tab.inner_text().strip()[:50]
            if text:
                analysis["tabs"].append(text)
    except Exception as e:
        print(f"Tabs 分析: {e}")

    # 分析工具栏
    try:
        toolbar = page.locator('.page-toolbar button, .toolbar button, .ant-card-extra button').all()
        for btn in toolbar[:20]:
            text = btn.inner_text().strip()
            if text:
                classes = btn.get_attribute("class") or ""
                analysis["toolbar_buttons"].append({
                    "text": text,
                    "is_primary": "primary" in classes
                })
    except Exception as e:
        print(f"工具栏分析: {e}")

    # 分析搜索
    try:
        inputs = page.locator('.ant-input, .vxe-input--inner').all()
        for inp in inputs[:10]:
            placeholder = inp.get_attribute("placeholder") or ""
            if placeholder:
                analysis["search_inputs"].append(placeholder)
    except Exception as e:
        print(f"搜索分析: {e}")

    # 分析表格
    try:
        headers = page.locator('.vxe-header--column .vxe-cell--title, .ant-table-thead th').all()
        for h in headers[:30]:
            text = h.inner_text().strip()
            if text and len(text) < 30:
                analysis["table_columns"].append(text)
    except Exception as e:
        print(f"表格分析: {e}")

    # 尝试点击新建按钮
    try:
        create_btns = page.locator('button:has-text("新建"), button:has-text("申请"), button:has-text("创建")').all()
        if create_btns:
            create_btns[0].click()
            page.wait_for_timeout(3000)

            # 截图弹窗
            dialog_path = os.path.join(OUTPUT_DIR, page_info["name"], "dialog.png")
            page.screenshot(path=dialog_path)

            # 分析弹窗字段
            dialog = page.locator('.ant-modal-content, [role="dialog"]').first
            if dialog.is_visible():
                title = dialog.locator('.ant-modal-title').inner_text().strip()

                form_labels = dialog.locator('.ant-form-item-label label').all()
                fields = []
                for label in form_labels:
                    text = label.inner_text().strip()
                    if text:
                        fields.append(text)

                analysis["create_dialog"] = {
                    "title": title,
                    "fields": fields
                }

                # 关闭弹窗
                cancel_btn = page.locator('.ant-modal-close, button:has-text("取消")')
                if cancel_btn.count() > 0:
                    cancel_btn.first.click()
                    page.wait_for_timeout(1000)
    except Exception as e:
        print(f"新建弹窗分析: {e}")

    # 保存分析
    result_path = os.path.join(OUTPUT_DIR, page_info["name"], "result.json")
    with open(result_path, 'w', encoding='utf-8') as f:
        json.dump(analysis, f, ensure_ascii=False, indent=2)
    print(f"结果: {result_path}")

    return analysis

def main():
    setup()

    with sync_playwright() as p:
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--ignore-certificate-errors-spki-list']
        )
        context = browser.new_context(ignore_https_errors=True)
        page = context.new_page()

        # 登录
        login(page)

        # 捕获 API
        api_calls = capture_api_calls(page)

        # 分析页面
        results = []
        for page_info in PAGES:
            result = analyze_page(page, page_info, api_calls)
            results.append(result)

        browser.close()

        # 打印汇总
        print("\n" + "="*50)
        print("分析汇总")
        print("="*50)
        for r in results:
            print(f"\n【{r['page']}】")
            print(f"  URL: {r['url']}")
            print(f"  API: {len(r['api_calls'])} 个")
            print(f"  工具栏: {r['toolbar_buttons']}")
            print(f"  表格列: {r['table_columns']}")
            if r['create_dialog'].get('fields'):
                print(f"  新建字段: {r['create_dialog']['fields']}")

if __name__ == "__main__":
    main()