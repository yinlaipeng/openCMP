#!/usr/bin/env python3
"""
Cloudpods 存储模块页面分析脚本
分析块存储、存储桶、表格存储、文件系统、NAS权限组页面设计
"""

from playwright.sync_api import sync_playwright
import json
import os
import time

# 配置
CLOUDPODS_BASE_URL = "https://127.0.0.1"
USERNAME = "admin"
PASSWORD = "admin@123"
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/cloudpods_analysis/output_storage"

# 要分析的页面
PAGES_TO_ANALYZE = {
    "block_storage": "/blockstorage",
    "buckets": "/bucket",
    "table_storage": "/table-storage",
    "file_systems": "/nas",
    "nas_groups": "/access-group"
}

def ensure_output_dir():
    """确保输出目录存在"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def login_to_cloudpods(page):
    """登录到 Cloudpods 系统"""
    print("正在登录 Cloudpods...")

    # 导航到登录页面
    page.goto(CLOUDPODS_BASE_URL, wait_until="networkidle")
    page.wait_for_timeout(2000)

    # 检查是否已经在登录页面或需要登录
    current_url = page.url
    print(f"当前 URL: {current_url}")

    # 尝试查找登录表单
    try:
        # 查找用户名输入框
        username_input = page.locator('input[type="text"], input[name="username"], input[placeholder*="用户"], input[placeholder*="Username"]').first
        if username_input.is_visible():
            username_input.fill(USERNAME)
            print("已输入用户名")

            # 查找密码输入框
            password_input = page.locator('input[type="password"], input[name="password"]').first
            password_input.fill(PASSWORD)
            print("已输入密码")

            # 查找登录按钮
            login_button = page.locator('button:has-text("登录"), button:has-text("Login"), button[type="submit"]').first
            login_button.click()
            print("已点击登录按钮")

            # 等待登录完成
            page.wait_for_load_state("networkidle")
            page.wait_for_timeout(3000)
            print(f"登录后 URL: {page.url}")
    except Exception as e:
        print(f"登录过程: {e}")

    return page

def analyze_page(page, page_name, page_path):
    """分析单个页面的设计"""
    print(f"\n分析页面: {page_name} ({page_path})")

    full_url = f"{CLOUDPODS_BASE_URL}{page_path}"
    page.goto(full_url, wait_until="networkidle")
    page.wait_for_timeout(2000)

    # 截取页面截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{page_name}_screenshot.png")
    page.screenshot(path=screenshot_path, full_page=True)
    print(f"截图已保存: {screenshot_path}")

    # 获取页面 HTML 结构
    html_content = page.content()
    html_path = os.path.join(OUTPUT_DIR, f"{page_name}_html.html")
    with open(html_path, 'w', encoding='utf-8') as f:
        f.write(html_content)
    print(f"HTML 已保存: {html_path}")

    # 分析页面元素
    analysis_result = {
        "page_name": page_name,
        "page_path": page_path,
        "url": full_url,
        "title": page.title(),
        "elements": {}
    }

    # 查找主要 UI 元素
    try:
        # 查找卡片/容器
        cards = page.locator('.el-card, .card, [class*="card"]').all()
        analysis_result["elements"]["cards"] = len(cards)

        # 查找表格
        tables = page.locator('.el-table, table, [class*="table"]').all()
        analysis_result["elements"]["tables"] = len(tables)

        # 查找按钮
        buttons = page.locator('.el-button, button, [class*="btn"]').all()
        button_texts = []
        for btn in buttons[:20]:  # 只取前20个
            try:
                text = btn.inner_text()
                if text.strip():
                    button_texts.append(text.strip())
            except:
                pass
        analysis_result["elements"]["buttons"] = button_texts

        # 查找搜索表单/筛选器
        search_forms = page.locator('.el-form, .search-form, form, [class*="search"]').all()
        analysis_result["elements"]["search_forms"] = len(search_forms)

        # 查找输入框
        inputs = page.locator('.el-input input, input[type="text"], input[type="search"]').all()
        input_count = len(inputs)
        analysis_result["elements"]["inputs"] = input_count

        # 查找下拉选择框
        selects = page.locator('.el-select, select, [class*="select"]').all()
        analysis_result["elements"]["selects"] = len(selects)

        # 查找分页组件
        pagination = page.locator('.el-pagination, [class*="pagination"]').all()
        analysis_result["elements"]["pagination"] = len(pagination)

        # 查找标签页
        tabs = page.locator('.el-tabs, [class*="tabs"]').all()
        analysis_result["elements"]["tabs"] = len(tabs)

        # 获取标签页标签名称
        tab_labels = []
        try:
            tab_items = page.locator('.el-tabs__item').all()
            for tab in tab_items:
                try:
                    text = tab.inner_text()
                    if text.strip():
                        tab_labels.append(text.strip())
                except:
                    pass
        except:
            pass
        analysis_result["elements"]["tab_labels"] = tab_labels

        # 查找对话框/弹窗
        dialogs = page.locator('.el-dialog, [class*="modal"], [class*="dialog"]').all()
        analysis_result["elements"]["dialogs"] = len(dialogs)

        # 查找状态标签
        tags = page.locator('.el-tag, [class*="tag"]').all()
        tag_texts = []
        for tag in tags[:15]:
            try:
                text = tag.inner_text()
                if text.strip():
                    tag_texts.append(text.strip())
            except:
                pass
        analysis_result["elements"]["tags"] = tag_texts

        # 获取表格列信息
        table_headers = page.locator('.el-table__header th, table th, [class*="table-header"]').all()
        header_texts = []
        for header in table_headers:
            try:
                text = header.inner_text()
                if text.strip():
                    header_texts.append(text.strip())
            except:
                pass
        analysis_result["elements"]["table_headers"] = header_texts

        # 查找操作列按钮
        action_buttons = page.locator('.el-table [class*="operation"] button, .el-table .el-button--small').all()
        action_texts = []
        for btn in action_buttons[:10]:
            try:
                text = btn.inner_text()
                if text.strip():
                    action_texts.append(text.strip())
            except:
                pass
        analysis_result["elements"]["action_buttons"] = action_texts

    except Exception as e:
        print(f"元素分析错误: {e}")
        analysis_result["analysis_error"] = str(e)

    # 保存分析结果
    result_path = os.path.join(OUTPUT_DIR, f"{page_name}_analysis.json")
    with open(result_path, 'w', encoding='utf-8') as f:
        json.dump(analysis_result, f, ensure_ascii=False, indent=2)
    print(f"分析结果已保存: {result_path}")

    return analysis_result

def main():
    """主函数"""
    ensure_output_dir()

    print("=" * 60)
    print("Cloudpods 存储模块页面分析脚本")
    print("=" * 60)

    all_results = {}

    with sync_playwright() as p:
        # 启动浏览器，忽略 HTTPS 错误
        browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--ignore-certificate-errors-spki-list']
        )

        context = browser.new_context(
            ignore_https_errors=True,
            viewport={"width": 1920, "height": 1080}
        )

        page = context.new_page()

        # 登录
        login_to_cloudpods(page)

        # 分析每个页面
        for page_name, page_path in PAGES_TO_ANALYZE.items():
            try:
                result = analyze_page(page, page_name, page_path)
                all_results[page_name] = result
            except Exception as e:
                print(f"分析页面 {page_name} 时出错: {e}")
                all_results[page_name] = {"error": str(e)}

        browser.close()

    # 保存汇总结果
    summary_path = os.path.join(OUTPUT_DIR, "all_storage_pages_summary.json")
    with open(summary_path, 'w', encoding='utf-8') as f:
        json.dump(all_results, f, ensure_ascii=False, indent=2)
    print(f"\n汇总结果已保存: {summary_path}")

    print("\n" + "=" * 60)
    print("分析完成!")
    print("=" * 60)

if __name__ == "__main__":
    main()