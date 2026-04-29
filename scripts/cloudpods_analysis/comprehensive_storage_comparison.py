#!/usr/bin/env python3
"""
Cloudpods vs openCMP 存储模块综合对比验证脚本
使用 Playwright 同时分析两个系统的存储页面，进行详细对比
"""

from playwright.sync_api import sync_playwright
import json
import os

# 配置
CLOUDPODS_BASE_URL = "https://127.0.0.1"
OPENCMP_BASE_URL = "http://localhost:3003"  # 前端实际运行端口
USERNAME = "admin"
PASSWORD = "admin@123"
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/cloudpods_analysis/output_storage_verification"

# 要对比的页面映射
PAGES_TO_COMPARE = {
    "block_storage": {
        "cloudpods": "/blockstorage",
        "opencmp": "/storage/block/block-storage"
    },
    "buckets": {
        "cloudpods": "/bucket",
        "opencmp": "/storage/object/buckets"
    },
    "table_storage": {
        "cloudpods": "/table-storage",
        "opencmp": "/storage/table/table-storage"
    },
    "file_systems": {
        "cloudpods": "/nas",
        "opencmp": "/storage/file/file-systems"
    },
    "nas_groups": {
        "cloudpods": "/access-group",
        "opencmp": "/storage/file/nas-groups"
    }
}

def ensure_output_dir():
    """确保输出目录存在"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def login_cloudpods(page):
    """登录到 Cloudpods 系统"""
    print("正在登录 Cloudpods...")
    page.goto(CLOUDPODS_BASE_URL, wait_until="networkidle")
    page.wait_for_timeout(2000)

    try:
        username_input = page.locator('input[type="text"], input[name="username"], input[placeholder*="用户"], input[placeholder*="Username"]').first
        if username_input.is_visible():
            username_input.fill(USERNAME)
            password_input = page.locator('input[type="password"], input[name="password"]').first
            password_input.fill(PASSWORD)
            login_button = page.locator('button:has-text("登录"), button:has-text("Login"), button[type="submit"]').first
            login_button.click()
            page.wait_for_load_state("networkidle")
            page.wait_for_timeout(3000)
            print("Cloudpods 登录成功")
    except Exception as e:
        print(f"Cloudpods 登录: {e}")

def login_opencmp(page):
    """登录到 openCMP 系统"""
    print("正在登录 openCMP...")
    page.goto(OPENCMP_BASE_URL, wait_until="networkidle")
    page.wait_for_timeout(2000)

    try:
        # 检查是否需要登录
        username_input = page.locator('input[type="text"], input[name="username"], input[placeholder*="用户"]').first
        if username_input.is_visible():
            username_input.fill(USERNAME)
            password_input = page.locator('input[type="password"], input[name="password"]').first
            password_input.fill(PASSWORD)
            login_button = page.locator('button:has-text("登录"), button[type="submit"]').first
            login_button.click()
            page.wait_for_load_state("networkidle")
            page.wait_for_timeout(3000)
            print("openCMP 登录成功")
    except Exception as e:
        print(f"openCMP 登录: {e}")

def analyze_page_detailed(page, page_name, system_name, url):
    """详细分析单个页面"""
    print(f"\n分析 {system_name} 页面: {page_name} ({url})")

    page.goto(url, wait_until="networkidle")
    page.wait_for_timeout(2000)

    # 截图
    screenshot_path = os.path.join(OUTPUT_DIR, f"{page_name}_{system_name}_screenshot.png")
    page.screenshot(path=screenshot_path, full_page=True)

    # 详细分析
    analysis = {
        "page_name": page_name,
        "system": system_name,
        "url": url,
        "title": page.title(),
        "elements": {}
    }

    try:
        # 工具栏按钮
        toolbar = page.locator('.card-header .toolbar, .el-card__header .toolbar').first
        if toolbar.is_visible():
            buttons = toolbar.locator('.el-button, button').all()
            button_texts = []
            for btn in buttons:
                try:
                    text = btn.inner_text().strip()
                    if text:
                        button_texts.append(text)
                except:
                    pass
            analysis["elements"]["toolbar_buttons"] = button_texts

        # 搜索表单
        search_form = page.locator('.search-form, .el-form--inline').first
        if search_form.is_visible():
            inputs = search_form.locator('.el-input input, input').all()
            selects = search_form.locator('.el-select').all()
            analysis["elements"]["search_inputs"] = len(inputs)
            analysis["elements"]["search_selects"] = len(selects)

        # 表格表头
        table_headers = page.locator('.el-table__header th').all()
        header_texts = []
        for header in table_headers:
            try:
                text = header.inner_text().strip()
                if text:
                    header_texts.append(text)
            except:
                pass
        analysis["elements"]["table_headers"] = header_texts

        # 表格选择列
        selection_column = page.locator('.el-table__body .el-checkbox, .el-table .el-table-column--selection').all()
        analysis["elements"]["has_selection"] = len(selection_column) > 0

        # 分页组件
        pagination = page.locator('.el-pagination').all()
        analysis["elements"]["has_pagination"] = len(pagination) > 0

        # 标签页
        tabs = page.locator('.el-tabs__item').all()
        tab_labels = []
        for tab in tabs:
            try:
                text = tab.inner_text().strip()
                if text:
                    tab_labels.append(text)
            except:
                pass
        analysis["elements"]["tabs"] = tab_labels

        # 对话框
        dialogs = page.locator('.el-dialog').all()
        analysis["elements"]["visible_dialogs"] = len(dialogs)

    except Exception as e:
        analysis["analysis_error"] = str(e)

    return analysis

def compare_pages(cloudpods_result, opencmp_result):
    """对比两个系统的页面"""
    comparison = {
        "page_name": cloudpods_result["page_name"],
        "comparison_result": {},
        "matches": [],
        "differences": []
    }

    # 对比工具栏按钮
    cp_buttons = cloudpods_result["elements"].get("toolbar_buttons", [])
    oc_buttons = opencmp_result["elements"].get("toolbar_buttons", [])

    # 按钮对比（考虑翻译）
    button_mapping = {
        "View": "查看",
        "Create": "创建",
        "Batch Action": "批量操作",
        "Batch operations": "批量操作",
        "Tags": "标签",
        "Delete": "删除"
    }

    cp_buttons_cn = [button_mapping.get(b, b) for b in cp_buttons]
    if set(cp_buttons_cn) == set(oc_buttons):
        comparison["matches"].append("工具栏按钮一致")
        comparison["comparison_result"]["toolbar_buttons"] = "✅ 一致"
    else:
        comparison["differences"].append(f"工具栏按钮不一致: Cloudpods={cp_buttons_cn}, openCMP={oc_buttons}")
        comparison["comparison_result"]["toolbar_buttons"] = "❌ 不一致"

    # 对比表头
    cp_headers = cloudpods_result["elements"].get("table_headers", [])
    oc_headers = opencmp_result["elements"].get("table_headers", [])

    # 表头对比（考虑翻译）
    header_mapping = {
        "Name": "名称",
        "Status": "状态",
        "Enable Status": "启用状态",
        "Physical capacity": "物理容量",
        "Virtual capacity": "虚拟容量",
        "Platform": "平台",
        "Owner Domain": "所属域",
        "Region": "区域",
        "Operations": "操作",
        "Tags": "标签",
        "Permission": "读写权限",
        "Project": "项目",
        "Created At": "创建时间",
        "Cloud account": "云账号",
        "FileSystem Type": "文件系统类型",
        "Storage Type": "存储类型",
        "Protocol Type": "协议类型",
        "Billing Type": "计费方式",
        "Mount Target Count": "挂载点数量"
    }

    cp_headers_cn = [header_mapping.get(h, h) for h in cp_headers]
    if len(cp_headers_cn) > 0 and len(oc_headers) > 0:
        # 检查是否所有 Cloudpods 表头都在 openCMP 中
        missing_headers = [h for h in cp_headers_cn if h not in oc_headers]
        if not missing_headers:
            comparison["matches"].append("表头一致")
            comparison["comparison_result"]["table_headers"] = "✅ 一致"
        else:
            comparison["differences"].append(f"表头缺失: {missing_headers}")
            comparison["comparison_result"]["table_headers"] = "❌ 不一致"
    else:
        comparison["comparison_result"]["table_headers"] = "⚠ 无法对比"

    # 对比搜索表单
    cp_search = cloudpods_result["elements"].get("search_inputs", 0) + cloudpods_result["elements"].get("search_selects", 0)
    oc_search = opencmp_result["elements"].get("search_inputs", 0) + opencmp_result["elements"].get("search_selects", 0)

    if cp_search > 0 and oc_search > 0:
        comparison["matches"].append("搜索表单存在")
        comparison["comparison_result"]["search_form"] = "✅ 一致"
    elif cp_search > 0 and oc_search == 0:
        comparison["differences"].append("openCMP缺少搜索表单")
        comparison["comparison_result"]["search_form"] = "❌ 不一致"
    else:
        comparison["comparison_result"]["search_form"] = "✅ 一致"

    # 对比分页
    cp_pagination = cloudpods_result["elements"].get("has_pagination", False)
    oc_pagination = opencmp_result["elements"].get("has_pagination", False)

    if cp_pagination and oc_pagination:
        comparison["matches"].append("分页组件存在")
        comparison["comparison_result"]["pagination"] = "✅ 一致"
    elif not cp_pagination and not oc_pagination:
        comparison["comparison_result"]["pagination"] = "✅ 一致 (均无)"
    else:
        comparison["differences"].append("分页组件不一致")
        comparison["comparison_result"]["pagination"] = "❌ 不一致"

    # 对比选择列
    cp_selection = cloudpods_result["elements"].get("has_selection", False)
    oc_selection = opencmp_result["elements"].get("has_selection", False)

    if cp_selection and oc_selection:
        comparison["matches"].append("表格选择列存在")
        comparison["comparison_result"]["selection"] = "✅ 一致"
    elif not cp_selection and not oc_selection:
        comparison["comparison_result"]["selection"] = "✅ 一致 (均无)"
    else:
        comparison["differences"].append("表格选择列不一致")
        comparison["comparison_result"]["selection"] = "❌ 不一致"

    # 最终判断
    if len(comparison["differences"]) == 0:
        comparison["final_result"] = "✅ 完全一致"
    else:
        comparison["final_result"] = "❌ 存在差异"

    return comparison

def main():
    """主函数"""
    ensure_output_dir()

    print("=" * 70)
    print("Cloudpods vs openCMP 存储模块综合对比验证")
    print("=" * 70)

    all_results = {}

    with sync_playwright() as p:
        # Cloudpods 浏览器
        cp_browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors', '--ignore-certificate-errors-spki-list']
        )
        cp_context = cp_browser.new_context(
            ignore_https_errors=True,
            viewport={"width": 1920, "height": 1080}
        )
        cp_page = cp_context.new_page()

        # openCMP 浏览器
        oc_browser = p.chromium.launch(headless=True)
        oc_context = oc_browser.new_context(viewport={"width": 1920, "height": 1080})
        oc_page = oc_context.new_page()

        # 登录
        login_cloudpods(cp_page)
        login_opencmp(oc_page)

        # 对比每个页面
        for page_name, urls in PAGES_TO_COMPARE.items():
            cloudpods_url = f"{CLOUDPODS_BASE_URL}{urls['cloudpods']}"
            opencmp_url = f"{OPENCMP_BASE_URL}{urls['opencmp']}"

            # 分析 Cloudpods
            cp_result = analyze_page_detailed(cp_page, page_name, "cloudpods", cloudpods_url)

            # 分析 openCMP
            oc_result = analyze_page_detailed(oc_page, page_name, "opencmp", opencmp_url)

            # 对比
            comparison = compare_pages(cp_result, oc_result)

            all_results[page_name] = {
                "cloudpods": cp_result,
                "opencmp": oc_result,
                "comparison": comparison
            }

            print(f"\n{page_name} 对比结果: {comparison['final_result']}")

        cp_browser.close()
        oc_browser.close()

    # 保存结果
    result_path = os.path.join(OUTPUT_DIR, "comprehensive_comparison.json")
    with open(result_path, 'w', encoding='utf-8') as f:
        json.dump(all_results, f, ensure_ascii=False, indent=2)

    # 生成汇总报告
    summary = {
        "total_pages": len(PAGES_TO_COMPARE),
        "consistent_pages": 0,
        "inconsistent_pages": 0,
        "details": []
    }

    for page_name, result in all_results.items():
        comparison = result["comparison"]
        if "✅" in comparison["final_result"]:
            summary["consistent_pages"] += 1
        else:
            summary["inconsistent_pages"] += 1
        summary["details"].append({
            "page": page_name,
            "result": comparison["final_result"],
            "matches": comparison["matches"],
            "differences": comparison["differences"]
        })

    summary_path = os.path.join(OUTPUT_DIR, "comparison_summary.json")
    with open(summary_path, 'w', encoding='utf-8') as f:
        json.dump(summary, f, ensure_ascii=False, indent=2)

    print("\n" + "=" * 70)
    print("对比验证完成!")
    print(f"一致页面: {summary['consistent_pages']}/{summary['total_pages']}")
    print(f"不一致页面: {summary['inconsistent_pages']}/{summary['total_pages']}")
    print("=" * 70)

    # 打印详细结果
    for detail in summary["details"]:
        print(f"\n{detail['page']}: {detail['result']}")
        if detail['matches']:
            print(f"  匹配项: {', '.join(detail['matches'])}")
        if detail['differences']:
            print(f"  差异项: {', '.join(detail['differences'])}")

if __name__ == "__main__":
    main()