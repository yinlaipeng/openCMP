#!/usr/bin/env python3
"""
分析 Cloudpods WAF策略和应用程序服务页面 - 增强版

使用 Playwright 登录 Cloudpods 系统，分析：
1. WAF策略页面 (/waf)
2. 应用程序服务页面 (/webapp)

增强登录流程，处理多种登录表单布局。
"""

import json
import os
import sys
import time
from playwright.sync_api import sync_playwright, expect

# 输出目录
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/.specify/cloudpods_analysis"

def setup_output_dir():
    """创建输出目录"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def login_cloudpods(page, base_url: str, username: str, password: str):
    """登录 Cloudpods 系统 - 增强版"""
    login_url = f"{base_url}/login"
    print(f"[INFO] 正在访问登录页面: {login_url}")

    page.goto(login_url, wait_until="networkidle", timeout=60000)
    time.sleep(3)

    # 截取登录页面截图
    login_screenshot = os.path.join(OUTPUT_DIR, "login_page.png")
    page.screenshot(path=login_screenshot, full_page=True)
    print(f"[INFO] 登录页面截图: {login_screenshot}")

    # 打印页面信息
    print(f"[INFO] 当前URL: {page.url}")
    print(f"[INFO] 页面标题: {page.title()}")

    # 尝试多种登录表单元素定位
    try:
        # 方法1: Ant Design Vue 登录表单
        username_selectors = [
            'input#username',
            'input[name="username"]',
            'input.ant-input[type="text"]',
            '.ant-input-affix-wrapper input[type="text"]',
            'input[placeholder*="用户"]',
            'input[placeholder*="Username"]',
            'input[placeholder*="账号"]',
        ]

        password_selectors = [
            'input#password',
            'input[name="password"]',
            'input.ant-input[type="password"]',
            '.ant-input-affix-wrapper input[type="password"]',
            'input[placeholder*="密码"]',
            'input[placeholder*="Password"]',
        ]

        username_input = None
        password_input = None

        # 查找用户名输入框
        for selector in username_selectors:
            try:
                elem = page.locator(selector).first
                if elem.is_visible(timeout=1000):
                    username_input = elem
                    print(f"[INFO] 找到用户名输入框: {selector}")
                    break
            except:
                continue

        # 查找密码输入框
        for selector in password_selectors:
            try:
                elem = page.locator(selector).first
                if elem.is_visible(timeout=1000):
                    password_input = elem
                    print(f"[INFO] 找到密码输入框: {selector}")
                    break
            except:
                continue

        if username_input and password_input:
            print(f"[INFO] 正在输入登录凭证...")
            username_input.fill(username)
            time.sleep(0.5)
            password_input.fill(password)
            time.sleep(0.5)

            # 点击登录按钮
            login_button_selectors = [
                'button[type="submit"]',
                'button.ant-btn-primary',
                'button:has-text("登录")',
                'button:has-text("Login")',
                'button:has-text("登 录")',
                '.ant-btn.ant-btn-primary',
            ]

            login_button = None
            for selector in login_button_selectors:
                try:
                    elem = page.locator(selector).first
                    if elem.is_visible(timeout=1000):
                        login_button = elem
                        print(f"[INFO] 找到登录按钮: {selector}")
                        break
                except:
                    continue

            if login_button:
                login_button.click()
                print(f"[INFO] 点击登录按钮...")
                time.sleep(5)

                # 等待页面跳转
                page.wait_for_load_state("networkidle", timeout=30000)
                print(f"[INFO] 登录后URL: {page.url}")
                print(f"[INFO] 登录后标题: {page.title()}")

                # 截取登录后页面截图
                after_login_screenshot = os.path.join(OUTPUT_DIR, "after_login.png")
                page.screenshot(path=after_login_screenshot, full_page=True)
                print(f"[INFO] 登录后截图: {after_login_screenshot}")

                # 检查是否成功登录（URL不包含/login）
                if "/login" not in page.url:
                    print(f"[INFO] 登录成功!")
                    return True
                else:
                    print(f"[WARN] 可能登录失败，仍在登录页面")
                    return False

    except Exception as e:
        print(f"[ERROR] 登录流程异常: {e}")
        import traceback
        traceback.print_exc()
        return False

    print(f"[ERROR] 未找到登录表单元素")
    return False

def analyze_page_detailed(page, page_url: str, page_name: str):
    """详细分析单个页面的结构"""
    print(f"\n{'='*60}")
    print(f"[INFO] 正在分析页面: {page_name} ({page_url})")
    print(f"{'='*60}")

    result = {
        "page_name": page_name,
        "page_url": page_url,
        "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
        "current_url": page_url,
        "page_title": "",
        "structure": {},
        "toolbar": [],
        "search": {},
        "table": {},
        "api_calls": [],
        "modals": {}
    }

    # 导航到页面
    page.goto(page_url, wait_until="networkidle", timeout=60000)
    time.sleep(3)

    result["current_url"] = page.url

    # 截取全页截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{page_name}_full.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"[INFO] 截图保存: {screenshot_path}")

    # 获取页面标题
    result["page_title"] = page.title()
    print(f"[INFO] 页面标题: {result['page_title']}")

    # 检查是否在登录页面（说明没有成功登录或被重定向）
    if "/login" in page.url:
        print(f"[WARN] 被重定向到登录页面，无法分析目标页面")
        result["structure"]["redirected_to_login"] = True
        return result

    # 1. 分析页面结构 - 检查侧边栏菜单
    try:
        sidebar = page.locator('.level-2-wrap, .ant-layout-sider, aside, .sidebar').first
        if sidebar.is_visible():
            # 获取菜单项
            menu_items = sidebar.locator('.ant-menu-item, .menu-item, li').all()
            menu_texts = [item.text_content().strip() for item in menu_items if item.text_content().strip()]
            result["structure"]["sidebar_menu"] = menu_texts[:20]  # 只取前20项
            print(f"[INFO] 侧边栏菜单项: {len(menu_texts)} 个")
    except Exception as e:
        print(f"[WARN] 侧边栏分析异常: {e}")

    # 2. 分析页面标题区域
    try:
        page_header_selectors = [
            '.page-header',
            '.ant-page-header',
            'h2.page-title',
            'h3.page-title',
            '.card-header',
        ]
        for selector in page_header_selectors:
            elem = page.locator(selector).first
            if elem.is_visible(timeout=1000):
                header_text = elem.text_content().strip()
                result["structure"]["page_header"] = header_text
                print(f"[INFO] 页面标题区域: {header_text}")
                break
    except Exception as e:
        print(f"[WARN] 页面标题分析异常: {e}")

    # 3. 分析工具栏按钮
    toolbar_buttons = []
    try:
        toolbar_selectors = [
            '.page-toolbar',
            '.toolbar',
            '.ant-space-horizontal',
            '.card-header .buttons',
            '.top-bar',
        ]

        for selector in toolbar_selectors:
            toolbar = page.locator(selector).first
            if toolbar.is_visible(timeout=1000):
                buttons = toolbar.locator('button, .ant-btn').all()
                for btn in buttons:
                    try:
                        btn_text = btn.text_content().strip()
                        btn_class = btn.get_attribute("class") or ""
                        is_disabled = btn.is_disabled()
                        is_primary = "primary" in btn_class.lower()
                        toolbar_buttons.append({
                            "text": btn_text,
                            "type": btn_class,
                            "disabled": is_disabled,
                            "primary": is_primary
                        })
                    except:
                        continue
                break

        if not toolbar_buttons:
            # 尝试直接查找页面上的按钮
            all_buttons = page.locator('button.ant-btn, button.el-button').all()
            for btn in all_buttons[:10]:  # 只取前10个
                try:
                    btn_text = btn.text_content().strip()
                    if btn_text:
                        btn_class = btn.get_attribute("class") or ""
                        toolbar_buttons.append({
                            "text": btn_text,
                            "type": btn_class,
                            "disabled": btn.is_disabled()
                        })
                except:
                    continue

        print(f"[INFO] 工具栏按钮: {len(toolbar_buttons)} 个")
        for btn in toolbar_buttons:
            print(f"  - {btn['text']} ({btn['type'][:30]}) disabled={btn['disabled']}")
    except Exception as e:
        print(f"[WARN] 工具栏分析异常: {e}")
    result["toolbar"] = toolbar_buttons

    # 4. 分析搜索区域
    try:
        search_selectors = [
            '.search-box-wrap',
            '.search-bar',
            '.filter-card',
            '.ant-card.filter-card',
            '.search-area',
        ]

        for selector in search_selectors:
            search_area = page.locator(selector).first
            if search_area.is_visible(timeout=1000):
                # 搜索输入框
                try:
                    search_input = search_area.locator('input.ant-input, input.el-input__inner, input[type="text"]').first
                    if search_input.is_visible():
                        placeholder = search_input.get_attribute("placeholder") or ""
                        result["search"]["input_placeholder"] = placeholder
                        print(f"[INFO] 搜索框 placeholder: {placeholder}")
                except:
                    pass

                # 筛选下拉框
                try:
                    selects = search_area.locator('.ant-select, .el-select').all()
                    result["search"]["filter_count"] = len(selects)
                    print(f"[INFO] 筛选下拉框: {len(selects)} 个")
                except:
                    pass

                break
    except Exception as e:
        print(f"[WARN] 搜索区域分析异常: {e}")

    # 5. 分析表格结构
    try:
        table_selectors = [
            '.vxe-table',
            '.ant-table',
            '.el-table',
            'table',
        ]

        for selector in table_selectors:
            table = page.locator(selector).first
            if table.is_visible(timeout=1000):
                # 表头列
                try:
                    headers = table.locator('th, .vxe-header--column, .ant-table-thead th, .el-table__header th').all()
                    header_texts = []
                    for h in headers:
                        text = h.text_content().strip()
                        if text:
                            header_texts.append(text)
                    result["table"]["headers"] = header_texts
                    print(f"[INFO] 表格表头: {header_texts}")
                except:
                    pass

                # 表格行数
                try:
                    rows = table.locator('tbody tr, .vxe-body--row, .ant-table-tbody tr, .el-table__body tr').all()
                    result["table"]["row_count"] = len(rows)
                    print(f"[INFO] 表格行数: {len(rows)}")
                except:
                    pass

                # 分页
                try:
                    pagination = page.locator('.vxe-pager, .ant-pagination, .el-pagination').first
                    if pagination.is_visible():
                        result["table"]["has_pagination"] = True
                        # 获取分页信息
                        page_info = pagination.locator('.ant-pagination-total-text, .vxe-pager--total').first
                        if page_info.is_visible():
                            result["table"]["pagination_info"] = page_info.text_content().strip()
                            print(f"[INFO] 分页信息: {result['table']['pagination_info']}")
                except:
                    pass

                break
    except Exception as e:
        print(f"[WARN] 表格分析异常: {e}")

    # 6. 尝试打开新建弹窗分析
    try:
        create_btn_selectors = [
            'button.ant-btn-primary:has-text("新建")',
            'button.ant-btn-primary:has-text("Create")',
            'button.ant-btn-primary:has-text("添加")',
            'button:has-text("新建")',
            'button:has-text("Create")',
        ]

        create_btn = None
        for selector in create_btn_selectors:
            try:
                elem = page.locator(selector).first
                if elem.is_visible(timeout=1000) and not elem.is_disabled():
                    create_btn = elem
                    break
            except:
                continue

        if create_btn:
            print(f"[INFO] 找到新建按钮，尝试点击...")
            create_btn.click()
            time.sleep(2)

            # 截取弹窗截图
            modal_screenshot = os.path.join(OUTPUT_DIR, f"{page_name}_create_modal.png")
            page.screenshot(path=modal_screenshot)
            print(f"[INFO] 弹窗截图保存: {modal_screenshot}")

            # 分析弹窗结构
            modal_selectors = [
                '.ant-modal',
                '.el-dialog',
                '.vxe-modal',
            ]

            for selector in modal_selectors:
                modal = page.locator(selector).first
                if modal.is_visible(timeout=2000):
                    # 弹窗标题
                    try:
                        modal_title = modal.locator('.ant-modal-title, .el-dialog__title, .vxe-modal--title').first
                        if modal_title.is_visible():
                            result["modals"]["create_title"] = modal_title.text_content().strip()
                            print(f"[INFO] 弹窗标题: {result['modals']['create_title']}")
                    except:
                        pass

                    # 弹窗表单字段
                    try:
                        form_items = modal.locator('.ant-form-item, .el-form-item').all()
                        form_fields = []
                        for item in form_items:
                            try:
                                label = item.locator('.ant-form-item-label label, .el-form-item__label').first
                                if label.is_visible():
                                    field_label = label.text_content().strip()
                                    required = item.locator('.ant-form-item-required, .required').is_visible()
                                    form_fields.append({
                                        "label": field_label,
                                        "required": required
                                    })
                            except:
                                continue
                        result["modals"]["create_fields"] = form_fields
                        print(f"[INFO] 弹窗表单字段: {len(form_fields)} 个")
                        for field in form_fields:
                            print(f"  - {field['label']} (required={field['required']})")
                    except:
                        pass

                    # 关闭弹窗
                    try:
                        close_btn = modal.locator('.ant-modal-close, .el-dialog__headerbtn, .vxe-modal--close-btn').first
                        if close_btn.is_visible():
                            close_btn.click()
                            time.sleep(1)
                    except:
                        # 使用 ESC 键关闭
                        page.keyboard.press("Escape")
                        time.sleep(1)

                    break
    except Exception as e:
        print(f"[WARN] 弹窗分析异常: {e}")

    # 7. 检查页面上的 Tabs
    try:
        tabs = page.locator('.ant-tabs-tab, .el-tabs__item').all()
        if tabs:
            tab_texts = [tab.text_content().strip() for tab in tabs if tab.text_content().strip()]
            result["structure"]["tabs"] = tab_texts
            print(f"[INFO] 页面Tabs: {tab_texts}")
    except Exception as e:
        print(f"[WARN] Tabs分析异常: {e}")

    return result

def capture_network_requests(page):
    """捕获网络请求"""
    api_calls = []

    def handle_request(request):
        url = request.url
        method = request.method
        # 只记录 API 调用
        if "/api/" in url:
            api_calls.append({
                "url": url,
                "method": method,
                "resource_type": request.resource_type
            })

    page.on("request", handle_request)
    return api_calls

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

    results = {
        "login_result": {},
        "pages": []
    }

    with sync_playwright() as p:
        # 启动浏览器（忽略 SSL 证书）
        browser = p.chromium.launch(
            headless=True,
            args=[
                '--ignore-certificate-errors',
                '--ignore-certificate-errors-spki-list',
                '--disable-web-security'
            ]
        )

        # 创建浏览器上下文（忽略 SSL）
        context = browser.new_context(
            ignore_https_errors=True,
            viewport={"width": 1920, "height": 1080}
        )

        page = context.new_page()

        # 设置请求捕获
        all_api_calls = capture_network_requests(page)

        try:
            # 登录
            login_success = login_cloudpods(page, base_url, username, password)
            results["login_result"] = {
                "success": login_success,
                "timestamp": time.strftime("%Y-%m-%d %H:%M:%S")
            }

            if not login_success:
                print(f"[WARN] 登录失败，尝试继续分析（可能会被重定向）...")

            # 分析各个页面
            for page_info in pages_to_analyze:
                page_url = f"{base_url}{page_info['url']}"
                result = analyze_page_detailed(page, page_url, page_info['name'])
                results["pages"].append(result)

            # 添加捕获的 API 调用
            results["captured_api_calls"] = all_api_calls[:50]  # 只取前50个

            # 保存分析结果
            output_file = os.path.join(OUTPUT_DIR, "analysis_results.json")
            with open(output_file, 'w', encoding='utf-8') as f:
                json.dump(results, f, ensure_ascii=False, indent=2)
            print(f"\n[INFO] 分析结果保存: {output_file}")

            # 打印 API 调用摘要
            print(f"\n[INFO] 捕获的 API 调用: {len(all_api_calls)} 个")
            waf_apis = [api for api in all_api_calls if "waf" in api["url"].lower()]
            webapp_apis = [api for api in all_api_calls if "webapp" in api["url"].lower()]
            print(f"[INFO] WAF相关API: {len(waf_apis)} 个")
            print(f"[INFO] Webapp相关API: {len(webapp_apis)} 个")

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