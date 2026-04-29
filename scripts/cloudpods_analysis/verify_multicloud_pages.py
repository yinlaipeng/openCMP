#!/usr/bin/env python3
"""
多云管理模块对比验证脚本
对比 Cloudpods 和 openCMP 的多云管理页面设计
"""

from playwright.sync_api import sync_playwright
import json
import os

# Cloudpods 配置
CLOUDPODS_BASE_URL = "https://127.0.0.1"
CLOUDPODS_USERNAME = "admin"
CLOUDPODS_PASSWORD = "admin@123"

# openCMP 配置
OPENCMP_BASE_URL = "http://localhost:3003"

# 输出目录
OUTPUT_DIR = "/Users/aurora/Desktop/xtwork/git/openCMP/scripts/cloudpods_analysis/output_multicloud_verification"

# 要对比的页面
PAGES_TO_COMPARE = {
    "cloudaccount": {
        "cloudpods": "/cloudaccount",
        "opencmp": "/cloud-accounts"
    },
    "cloudgroup": {
        "cloudpods": "/cloudgroup",
        "opencmp": "/cloud-management/cloud-user-groups"
    },
    "proxysetting": {
        "cloudpods": "/proxysetting",
        "opencmp": "/cloud-management/proxies"
    },
    "projectmapping": {
        "cloudpods": "/projectmapping",
        "opencmp": "/cloud-management/sync-policies"
    }
}

def ensure_output_dir():
    """确保输出目录存在"""
    os.makedirs(OUTPUT_DIR, exist_ok=True)

def login_to_cloudpods(page):
    """登录到 Cloudpods 系统"""
    print("正在登录 Cloudpods...")
    page.goto(CLOUDPODS_BASE_URL, wait_until="networkidle")
    page.wait_for_timeout(2000)

    try:
        username_input = page.locator('input[type="text"], input[name="username"]').first
        if username_input.is_visible():
            username_input.fill(CLOUDPODS_USERNAME)
            password_input = page.locator('input[type="password"]').first
            password_input.fill(CLOUDPODS_PASSWORD)
            login_button = page.locator('button:has-text("登录"), button[type="submit"]').first
            login_button.click()
            page.wait_for_load_state("networkidle")
            page.wait_for_timeout(3000)
            print(f"Cloudpods 登录成功，当前 URL: {page.url}")
    except Exception as e:
        print(f"Cloudpods 登录过程: {e}")

    return page

def login_to_opencmp(page):
    """登录到 openCMP 系统"""
    print("正在登录 openCMP...")
    page.goto(OPENCMP_BASE_URL, wait_until="networkidle")
    page.wait_for_timeout(2000)

    try:
        username_input = page.locator('input[type="text"], input[name="username"]').first
        if username_input.is_visible():
            username_input.fill("admin")
            password_input = page.locator('input[type="password"]').first
            password_input.fill("admin@123")
            login_button = page.locator('button:has-text("登录"), button[type="submit"]').first
            login_button.click()
            page.wait_for_load_state("networkidle")
            page.wait_for_timeout(3000)
            print(f"openCMP 登录成功，当前 URL: {page.url}")
    except Exception as e:
        print(f"openCMP 登录过程: {e}")

    return page

def analyze_page_toolbar(page, url, name):
    """分析页面工具栏按钮"""
    page.goto(url, wait_until="networkidle")
    page.wait_for_timeout(2000)

    # 查找工具栏按钮
    buttons = page.locator('.el-button, button').all()
    toolbar_buttons = []

    for btn in buttons[:15]:
        try:
            text = btn.inner_text()
            if text.strip() and len(text.strip()) < 20:
                toolbar_buttons.append(text.strip())
        except:
            pass

    # 查找表格表头
    headers = page.locator('.el-table__header th, table th').all()
    table_headers = []

    for header in headers[:15]:
        try:
            text = header.inner_text()
            if text.strip():
                table_headers.append(text.strip())
        except:
            pass

    # 检查表格是否存在（即使没有数据）
    table_exists = len(page.locator('.el-table').all()) > 0

    return {
        "name": name,
        "url": url,
        "toolbar_buttons": toolbar_buttons,
        "table_headers": table_headers,
        "table_exists": table_exists
    }

def main():
    """主函数"""
    ensure_output_dir()

    print("=" * 60)
    print("多云管理模块对比验证")
    print("=" * 60)

    comparison_results = {}

    with sync_playwright() as p:
        # 分析 Cloudpods
        print("\n--- 分析 Cloudpods ---")
        cloudpods_browser = p.chromium.launch(
            headless=True,
            args=['--ignore-certificate-errors']
        )
        cloudpods_context = cloudpods_browser.new_context(
            ignore_https_errors=True,
            viewport={"width": 1920, "height": 1080}
        )
        cloudpods_page = cloudpods_context.new_page()

        login_to_cloudpods(cloudpods_page)

        for page_name, paths in PAGES_TO_COMPARE.items():
            cloudpods_url = f"{CLOUDPODS_BASE_URL}{paths['cloudpods']}"
            print(f"分析 Cloudpods: {page_name} ({cloudpods_url})")

            try:
                result = analyze_page_toolbar(cloudpods_page, cloudpods_url, page_name)
                comparison_results[f"cloudpods_{page_name}"] = result
            except Exception as e:
                print(f"Cloudpods {page_name} 分析错误: {e}")
                comparison_results[f"cloudpods_{page_name}"] = {"error": str(e)}

        cloudpods_browser.close()

        # 分析 openCMP
        print("\n--- 分析 openCMP ---")
        opencmp_browser = p.chromium.launch(headless=True)
        opencmp_context = opencmp_browser.new_context(
            viewport={"width": 1920, "height": 1080}
        )
        opencmp_page = opencmp_context.new_page()

        login_to_opencmp(opencmp_page)

        for page_name, paths in PAGES_TO_COMPARE.items():
            opencmp_url = f"{OPENCMP_BASE_URL}{paths['opencmp']}"
            print(f"分析 openCMP: {page_name} ({opencmp_url})")

            try:
                result = analyze_page_toolbar(opencmp_page, opencmp_url, page_name)
                comparison_results[f"opencmp_{page_name}"] = result
            except Exception as e:
                print(f"openCMP {page_name} 分析错误: {e}")
                comparison_results[f"opencmp_{page_name}"] = {"error": str(e)}

        opencmp_browser.close()

    # 对比分析
    print("\n--- 对比结果 ---")
    comparison_summary = {}

    for page_name in PAGES_TO_COMPARE.keys():
        cloudpods_result = comparison_results.get(f"cloudpods_{page_name}", {})
        opencmp_result = comparison_results.get(f"opencmp_{page_name}", {})

        if "error" in cloudpods_result or "error" in opencmp_result:
            comparison_summary[page_name] = {
                "status": "error",
                "message": "分析出错"
            }
            continue

        # 检查表格是否存在
        opencmp_table_exists = opencmp_result.get("table_exists", False)
        opencmp_buttons = opencmp_result.get("toolbar_buttons", [])
        opencmp_headers = opencmp_result.get("table_headers", [])

        comparison_summary[page_name] = {
            "status": "checked",
            "cloudpods_buttons": cloudpods_result.get("toolbar_buttons", []),
            "opencmp_buttons": opencmp_buttons,
            "cloudpods_headers": cloudpods_result.get("table_headers", []),
            "opencmp_headers": opencmp_headers,
            "opencmp_table_exists": opencmp_table_exists
        }

        print(f"\n{page_name}:")
        print(f"  Cloudpods 按钮: {cloudpods_result.get('toolbar_buttons', [])}")
        print(f"  openCMP 按钮: {opencmp_buttons}")
        print(f"  openCMP 表头: {opencmp_headers}")
        print(f"  openCMP 表格存在: {opencmp_table_exists}")

    # 保存结果
    result_path = os.path.join(OUTPUT_DIR, "multicloud_comparison.json")
    with open(result_path, 'w', encoding='utf-8') as f:
        json.dump(comparison_results, f, ensure_ascii=False, indent=2)

    summary_path = os.path.join(OUTPUT_DIR, "comparison_summary.json")
    with open(summary_path, 'w', encoding='utf-8') as f:
        json.dump(comparison_summary, f, ensure_ascii=False, indent=2)

    print(f"\n结果已保存: {OUTPUT_DIR}")
    print("=" * 60)

if __name__ == "__main__":
    main()