#!/usr/bin/env python3
"""
Cloudpods 网络服务页面详细分析脚本 - 改进版
使用 Puppeteer MCP 工具风格分析
"""

import json
import os
import time
from playwright.sync_api import sync_playwright, TimeoutError as PlaywrightTimeout

# 配置
BASE_URL = "https://127.0.0.1"
USERNAME = "admin"
PASSWORD = "admin@123"

# 要分析的页面
PAGES_TO_ANALYZE = [
    {"name": "EIP弹性公网IP", "url": "/eip"},
    {"name": "NAT网关", "url": "/nat"},
    {"name": "DNS解析", "url": "/vpc-network"},
    {"name": "IPv6网关", "url": "/ipv6-gateway"},
]

OUTPUT_DIR = "/tmp/cloudpods_network_analysis"

def setup():
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    for p in PAGES_TO_ANALYZE:
        os.makedirs(os.path.join(OUTPUT_DIR, p["name"]), exist_ok=True)

def login_with_retry(page, max_retries=3):
    """带重试的登录"""
    for attempt in range(max_retries):
        print(f"登录尝试 {attempt + 1}/{max_retries}")

        try:
            # 导航到登录页
            page.goto(BASE_URL + "/login", wait_until="networkidle", timeout=30000)
            page.wait_for_timeout(2000)

            # 检查是否已登录（页面可能直接跳转到主页）
            if page.url.endswith("/login") or "/login" in page.url:
                # 使用多种选择器尝试找到用户名输入框
                username_selectors = [
                    'input[placeholder*="用户"]',
                    'input[placeholder*="Username"]',
                    'input[name="username"]',
                    '#username',
                    '.login-input:first input',
                    'input.ant-input:first-of-type',
                ]

                username_input = None
                for sel in username_selectors:
                    try:
                        if page.locator(sel).count() > 0:
                            username_input = page.locator(sel).first
                            if username_input.is_visible():
                                break
                    except:
                        continue

                if username_input:
                    username_input.fill(USERNAME)
                    page.wait_for_timeout(500)

                    # 填写密码
                    page.locator('input[type="password"]').fill(PASSWORD)
                    page.wait_for_timeout(500)

                    # 点击登录
                    login_btn_selectors = [
                        'button[type="submit"]',
                        'button.ant-btn-primary',
                        'button:has-text("登录")',
                        'button:has-text("Login")',
                    ]

                    for sel in login_btn_selectors:
                        try:
                            btn = page.locator(sel).first
                            if btn.is_visible():
                                btn.click()
                                break
                        except:
                            continue

                    # 等待登录完成
                    page.wait_for_load_state("networkidle", timeout=15000)
                    page.wait_for_timeout(3000)

                    # 检查是否成功
                    if not "/login" in page.url:
                        print("登录成功!")
                        return True
                else:
                    print("未找到用户名输入框")
            else:
                print("可能已登录，页面已跳转")
                return True

        except PlaywrightTimeout:
            print("超时，重试...")
        except Exception as e:
            print(f"登录错误: {e}")

        time.sleep(2)

    return False

def analyze_network_requests(page):
    """捕获网络请求"""
    requests = []

    def on_request(request):
        if "/api/" in request.url:
            requests.append({
                "url": request.url,
                "method": request.method,
                "resource_type": request.resource_type
            })

    page.on("request", on_request)
    return requests

def analyze_page(page, page_info, requests):
    """详细分析页面"""
    print(f"\n{'='*60}")
    print(f"分析: {page_info['name']} ({page_info['url']})")
    print(f"{'='*60}")

    url = BASE_URL + page_info["url"]

    # 导航到页面
    page.goto(url, wait_until="networkidle", timeout=30000)
    page.wait_for_timeout(3000)

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, page_info["name"], "page.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"截图: {screenshot_path}")

    # 分析页面结构
    analysis = {
        "page_name": page_info["name"],
        "page_url": page_info["url"],
        "full_url": page.url,
        "api_calls": [],
        "toolbar": {},
        "search": {},
        "table": {},
        "tabs": [],
        "create_dialog": {},
    }

    # 过滤相关 API
    page_url_path = page_info["url"].replace("/", "")
    for req in requests:
        if page_url_path in req["url"] or "/eip" in req["url"] or "/nat" in req["url"]:
            analysis["api_calls"].append(req)

    # 分析 Tabs
    try:
        tabs = page.locator('.ant-tabs-tab, .vxe-tabs--header .vxe-tabs--item').all()
        for tab in tabs:
            text = tab.inner_text().strip()
            if text:
                analysis["tabs"].append(text)
    except:
        pass

    # 分析工具栏
    try:
        toolbar_selectors = ['.page-toolbar', '.toolbar', '.ant-card-extra', '.vxe-toolbar']
        for sel in toolbar_selectors:
            toolbar = page.locator(sel)
            if toolbar.count() > 0:
                toolbar_el = toolbar.first
                buttons = toolbar_el.locator('button').all()
                btn_info = []
                for btn in buttons:
                    text = btn.inner_text().strip()
                    classes = btn.get_attribute("class") or ""
                    btn_info.append({
                        "text": text,
                        "is_primary": "primary" in classes,
                        "is_disabled": btn.is_disabled()
                    })
                if btn_info:
                    analysis["toolbar"] = {"buttons": btn_info}
                    break
    except Exception as e:
        print(f"工具栏分析失败: {e}")

    # 分析搜索区域
    try:
        search_selectors = ['.search-box-wrap', '.ant-card-body .ant-input', '.vxe-toolbar .vxe-input']
        for sel in search_selectors:
            search = page.locator(sel)
            if search.count() > 0:
                analysis["search"] = {"selector": sel, "count": search.count()}
                break
    except:
        pass

    # 分析表格
    try:
        table_selectors = ['.vxe-table', '.ant-table', 'table']
        for sel in table_selectors:
            table = page.locator(sel)
            if table.count() > 0:
                # 获取列头
                headers = table.first.locator('th, .vxe-header--column').all()
                columns = []
                for h in headers:
                    text = h.inner_text().strip()
                    if text and text not in ['', '\n']:
                        columns.append(text)

                # 获取数据行数
                rows = table.first.locator('tbody tr, .vxe-table--body tr').all()

                analysis["table"] = {
                    "columns": columns,
                    "row_count": len(rows)
                }
                break
    except Exception as e:
        print(f"表格分析失败: {e}")

    # 尝试打开新建弹窗
    try:
        create_selectors = [
            'button:has-text("新建")',
            'button:has-text("申请")',
            'button:has-text("创建")',
            '.ant-btn-primary:has-text("新")',
        ]

        for sel in create_selectors:
            create_btn = page.locator(sel)
            if create_btn.count() > 0 and create_btn.first.is_visible() and not create_btn.first.is_disabled():
                create_btn.first.click()
                page.wait_for_timeout(2000)

                # 截图弹窗
                dialog_path = os.path.join(OUTPUT_DIR, page_info["name"], "create_dialog.png")
                page.screenshot(path=dialog_path)

                # 分析弹窗
                dialog = page.locator('.ant-modal-content, .ant-modal').first
                if dialog.is_visible():
                    title = dialog.locator('.ant-modal-title, .ant-modal-header').inner_text().strip()

                    # 表单项
                    form_items = dialog.locator('.ant-form-item').all()
                    fields = []
                    for item in form_items:
                        label = item.locator('.ant-form-item-label label').inner_text().strip()
                        # 判断输入类型
                        input_el = item.locator('input').first
                        select_el = item.locator('.ant-select').first
                        textarea_el = item.locator('textarea').first

                        field_type = "text"
                        if select_el.is_visible():
                            field_type = "select"
                        elif textarea_el.is_visible():
                            field_type = "textarea"
                        elif input_el.is_visible():
                            ipt_type = input_el.get_attribute("type") or "text"
                            field_type = ipt_type

                        fields.append({"label": label, "type": field_type})

                    analysis["create_dialog"] = {
                        "title": title,
                        "fields": fields
                    }

                    # 关闭弹窗
                    close_btn = page.locator('.ant-modal-close, button:has-text("取消")').first
                    if close_btn.is_visible():
                        close_btn.click()
                        page.wait_for_timeout(1000)

                break
    except Exception as e:
        print(f"新建弹窗分析失败: {e}")

    # 保存分析结果
    result_path = os.path.join(OUTPUT_DIR, page_info["name"], "analysis.json")
    with open(result_path, 'w', encoding='utf-8') as f:
        json.dump(analysis, f, ensure_ascii=False, indent=2)

    print(f"分析结果: {result_path}")
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
        if not login_with_retry(page):
            print("登录失败，退出")
            browser.close()
            return

        # 捕获 API 请求
        requests = analyze_network_requests(page)

        # 分析每个页面
        all_results = []
        for page_info in PAGES_TO_ANALYZE:
            result = analyze_page(page, page_info, requests)
            all_results.append(result)

        # 保存 HAR
        har_path = os.path.join(OUTPUT_DIR, "requests.json")
        with open(har_path, 'w', encoding='utf-8') as f:
            json.dump(requests, f, ensure_ascii=False, indent=2)

        browser.close()

        # 打印汇总
        print("\n" + "="*60)
        print("分析汇总")
        print("="*60)

        for r in all_results:
            print(f"\n【{r['page_name']}】")
            print(f"  URL: {r['full_url']}")
            print(f"  API调用: {len(r['api_calls'])} 个")
            if r['tabs']:
                print(f"  Tabs: {r['tabs']}")
            if r['toolbar'].get('buttons'):
                print(f"  工具栏按钮: {[b['text'] for b in r['toolbar']['buttons']]}")
            if r['table'].get('columns'):
                print(f"  表格列: {r['table']['columns']}")
            if r['create_dialog'].get('fields'):
                print(f"  新建弹窗字段: {[f['label'] for f in r['create_dialog']['fields']]}")

if __name__ == "__main__":
    main()