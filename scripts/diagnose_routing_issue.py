#!/usr/bin/env python3
"""
诊断前端路由切换问题
问题：打开 /compute/vms 页面后点击其他页面，浏览器 URL 变了但页面内容一直是虚拟机页面
"""

from playwright.sync_api import sync_playwright
import json
import time

def diagnose_routing():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)  # 非headless便于观察
        context = browser.new_context()
        page = context.new_page()

        # 捕获控制台日志
        console_logs = []
        page.on("console", lambda msg: console_logs.append({
            "type": msg.type,
            "text": msg.text
        }))

        # 捕获网络请求
        requests_log = []
        page.on("request", lambda req: requests_log.append({
            "url": req.url,
            "method": req.method
        }))

        results = {
            "steps": [],
            "console_errors": [],
            "router_state": [],
            "page_content": []
        }

        try:
            # 1. 登录
            print("Step 1: 登录系统")
            page.goto('http://localhost:3000/login')
            page.wait_for_load_state('networkidle')

            page.fill('input[placeholder="请输入用户名"]', 'admin')
            page.fill('input[placeholder="请输入密码"]', 'admin@123')
            page.click('button:has-text("登录")')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            # 检查是否跳转到 Dashboard
            current_url = page.url
            print(f"登录后 URL: {current_url}")
            results["steps"].append({
                "action": "登录",
                "url": current_url,
                "success": "/dashboard" in current_url or "/" in current_url
            })

            # 2. 访问虚拟机页面
            print("\nStep 2: 访问 /compute/vms")
            page.goto('http://localhost:3000/compute/vms')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            vms_url = page.url
            vms_title = page.title()
            vms_content = page.content()

            # 检查页面是否包含虚拟机相关内容
            vms_indicators = [
                "虚拟机管理",
                "VM",
                "主机",
                "vms",
                "compute"
            ]
            vms_found = any(ind in vms_content for ind in vms_indicators)

            print(f"虚拟机页面 URL: {vms_url}")
            print(f"虚拟机页面标题: {vms_title}")
            print(f"检测到虚拟机相关内容: {vms_found}")

            results["steps"].append({
                "action": "访问虚拟机页面",
                "url": vms_url,
                "title": vms_title,
                "vms_content_found": vms_found
            })

            # 截图保存虚拟机页面状态
            page.screenshot(path='/tmp/routing_vms_page.png')

            # 3. 获取 Vue Router 状态（通过 JavaScript）
            print("\nStep 3: 检查 Vue Router 状态")
            router_info = page.evaluate('''() => {
                if (window.__VUE_APP__) {
                    const router = window.__VUE_APP__.config.globalProperties.$router;
                    return {
                        currentRoute: router.currentRoute.value ? {
                            path: router.currentRoute.value.path,
                            name: router.currentRoute.value.name,
                            fullPath: router.currentRoute.value.fullPath
                        } : null,
                        routes: router.options.routes ? router.options.routes.map(r => r.path) : []
                    };
                }
                return { error: "Vue app not found" };
            }''')

            print(f"Router 状态: {json.dumps(router_info, indent=2)}")
            results["router_state"].append({
                "phase": "访问虚拟机页面后",
                "state": router_info
            })

            # 4. 点击侧边栏其他菜单
            print("\nStep 4: 点击侧边栏菜单切换页面")

            # 获取侧边栏菜单项
            sidebar_items = page.locator('.el-menu-item, .el-sub-menu__title').all()
            print(f"找到 {len(sidebar_items)} 个菜单项")

            # 尝试点击"主机模版"菜单
            # 先展开计算资源子菜单（如果需要）
            compute_menu = page.locator('.el-sub-menu:has-text("计算资源")')
            if compute_menu.count() > 0:
                print("展开计算资源菜单")
                compute_menu.click()
                time.sleep(1)

            # 点击主机模版
            host_templates_menu = page.locator('.el-menu-item:has-text("主机模版")')
            if host_templates_menu.count() > 0:
                print("点击主机模版菜单")
                host_templates_menu.click()
                time.sleep(3)  # 等待路由切换

                new_url = page.url
                new_content = page.content()

                # 检查 URL 是否变化
                url_changed = new_url != vms_url
                print(f"URL 是否变化: {url_changed}")
                print(f"新 URL: {new_url}")

                # 检查页面内容是否变化
                host_templates_indicators = [
                    "主机模版",
                    "host-templates",
                    "HostTemplate",
                    "模版"
                ]
                host_templates_found = any(ind in new_content for ind in host_templates_indicators)
                vms_still_found = any(ind in new_content for ind in vms_indicators)

                print(f"检测到主机模版内容: {host_templates_found}")
                print(f"虚拟机内容仍然存在: {vms_still_found}")

                results["steps"].append({
                    "action": "点击主机模版",
                    "url_before": vms_url,
                    "url_after": new_url,
                    "url_changed": url_changed,
                    "host_templates_found": host_templates_found,
                    "vms_still_found": vms_still_found,
                    "issue_detected": url_changed and vms_still_found and not host_templates_found
                })

                # 截图保存点击后状态
                page.screenshot(path='/tmp/routing_after_click.png')

                # 再次检查 Router 状态
                router_info_after = page.evaluate('''() => {
                    if (window.__VUE_APP__) {
                        const router = window.__VUE_APP__.config.globalProperties.$router;
                        return {
                            currentRoute: router.currentRoute.value ? {
                                path: router.currentRoute.value.path,
                                name: router.currentRoute.value.name,
                                fullPath: router.currentRoute.value.fullPath
                            } : null
                        };
                    }
                    return { error: "Vue app not found" };
                }''')

                print(f"点击后 Router 状态: {json.dumps(router_info_after, indent=2)}")
                results["router_state"].append({
                    "phase": "点击主机模版后",
                    "state": router_info_after
                })

                # 检查 router-view 是否更新
                router_view_content = page.evaluate('''() => {
                    const routerViews = document.querySelectorAll('.main-content, [class*="router-view"], .el-main');
                    if (routerViews.length > 0) {
                        return {
                            count: routerViews.length,
                            innerHTML: routerViews[0].innerHTML.substring(0, 500)
                        };
                    }
                    return { error: "router-view not found" };
                }''')

                print(f"Router-view 内容: {json.dumps(router_view_content, indent=2)}")
                results["page_content"].append(router_view_content)

            # 5. 测试刷新后是否正常
            print("\nStep 5: 刷新页面测试")
            page.reload()
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            refreshed_url = page.url
            refreshed_content = page.content()

            host_templates_found_refresh = any(ind in refreshed_content for ind in host_templates_indicators)
            vms_found_refresh = any(ind in refreshed_content for ind in vms_indicators)

            print(f"刷新后 URL: {refreshed_url}")
            print(f"刷新后主机模版内容: {host_templates_found_refresh}")
            print(f"刷新后虚拟机内容: {vms_found_refresh}")

            results["steps"].append({
                "action": "刷新页面",
                "url": refreshed_url,
                "host_templates_found": host_templates_found_refresh,
                "vms_found": vms_found_refresh,
                "refresh_fixes_issue": host_templates_found_refresh and not vms_found_refresh
            })

            page.screenshot(path='/tmp/routing_after_refresh.png')

            # 6. 分析问题根因
            print("\n=== 问题分析 ===")
            if results["steps"]:
                click_step = results["steps"][3] if len(results["steps"]) > 3 else {}
                if click_step.get("issue_detected"):
                    print("问题确认: URL变化但页面内容未更新")
                    print("可能原因:")
                    print("1. Vue Router 的 router-view 未正确响应路由变化")
                    print("2. 组件使用了 keep-alive 导致缓存")
                    print("3. 组件未正确 watch $route")
                    print("4. Layout 组件结构问题")

            # 收集控制台错误
            results["console_errors"] = [log for log in console_logs if log["type"] == "error"]

            # 保存完整报告
            with open('/tmp/routing_diagnosis_report.json', 'w') as f:
                json.dump(results, f, indent=2, ensure_ascii=False)

            print(f"\n报告已保存到 /tmp/routing_diagnosis_report.json")
            print(f"截图保存: /tmp/routing_vms_page.png, /tmp/routing_after_click.png, /tmp/routing_after_refresh.png")

        except Exception as e:
            print(f"错误: {e}")
            results["error"] = str(e)

        finally:
            browser.close()

        return results

if __name__ == "__main__":
    results = diagnose_routing()