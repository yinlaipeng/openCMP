#!/usr/bin/env python3
"""
详细诊断路由切换问题 - 检查 Vue 组件状态
"""

from playwright.sync_api import sync_playwright
import time
import json

def diagnose():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        results = {
            "steps": [],
            "vue_state": [],
            "dom_state": []
        }

        try:
            # 1. 登录
            print("Step 1: 登录系统")
            page.goto('http://localhost:3000/login')
            page.wait_for_load_state('networkidle')

            page.fill('input[placeholder="请输入用户名"]', 'admin')
            page.fill('input[placeholder="请输入密码"]', 'admin@123')
            page.click('button:has-text("登")')
            page.wait_for_load_state('networkidle')
            time.sleep(2)

            # 检查 token
            token = page.evaluate('localStorage.getItem("token")')
            print(f"Token: {token[:20] if token else 'None'}...")

            if not token:
                print("登录失败，尝试其他登录方式")
                page.fill('input[type="text"]', 'admin')
                page.fill('input[type="password"]', 'admin@123')
                page.click('button')
                time.sleep(3)
                token = page.evaluate('localStorage.getItem("token")')
                print(f"Retry Token: {token[:20] if token else 'None'}...")

            results["steps"].append({
                "action": "登录",
                "url": page.url,
                "token_exists": bool(token)
            })

            if not token:
                print("无法登录，跳过后续测试")
                return results

            # 2. 访问虚拟机页面
            print("\nStep 2: 访问虚拟机页面")
            page.goto('http://localhost:3000/compute/vms')
            page.wait_for_load_state('networkidle')
            time.sleep(3)

            # 获取页面内容关键信息
            vms_page_content = page.evaluate('''() => {
                const mainContent = document.querySelector('.main-content');
                const computeLayout = document.querySelector('.compute-layout');
                const vmsContainer = document.querySelector('.vms-container');

                return {
                    mainContentHTML: mainContent ? mainContent.innerHTML.substring(0, 1000) : 'not found',
                    computeLayoutExists: !!computeLayout,
                    vmsContainerExists: !!vmsContainer,
                    titleText: document.querySelector('.page-header h2')?.innerText || 'no title'
                };
            }''')
            print(f"虚拟机页面内容: {json.dumps(vms_page_content, indent=2, ensure_ascii=False)}")
            results["dom_state"].append({
                "phase": "虚拟机页面",
                "state": vms_page_content
            })

            # 检查 Vue Router 当前路由
            router_state = page.evaluate('''() => {
                const app = document.querySelector('#app')?.__vue_app__;
                if (app) {
                    const router = app.config.globalProperties.$router;
                    return {
                        currentRoute: router.currentRoute.value?.path,
                        matchedCount: router.currentRoute.value?.matched?.length,
                        matchedPaths: router.currentRoute.value?.matched?.map(m => m.path)
                    };
                }
                return { error: 'Vue app not found' };
            }''')
            print(f"Router 状态: {json.dumps(router_state, indent=2, ensure_ascii=False)}")
            results["vue_state"].append({
                "phase": "虚拟机页面",
                "state": router_state
            })

            page.screenshot(path='/tmp/routing_vms_detailed.png')

            # 3. 导航到主机模版页面（使用 router.push 模拟）
            print("\nStep 3: 使用 JavaScript 导航到主机模版")
            page.evaluate('''() => {
                const app = document.querySelector('#app')?.__vue_app__;
                if (app) {
                    const router = app.config.globalProperties.$router;
                    router.push('/compute/host-templates');
                }
            }''')
            time.sleep(3)
            page.wait_for_load_state('networkidle')

            # 再次检查页面内容
            templates_page_content = page.evaluate('''() => {
                const mainContent = document.querySelector('.main-content');
                const computeLayout = document.querySelector('.compute-layout');
                const templatesContainer = document.querySelector('.host-templates-container');
                const pageTitle = document.querySelector('.page-header h2')?.innerText;

                // 检查是否有虚拟机页面残留内容
                const vmsElements = document.querySelectorAll('[class*="vms"]');
                const templatesElements = document.querySelectorAll('[class*="host-templates"]');

                return {
                    mainContentHTML: mainContent ? mainContent.innerHTML.substring(0, 1000) : 'not found',
                    computeLayoutExists: !!computeLayout,
                    templatesContainerExists: !!templatesContainer,
                    pageTitle: pageTitle || 'no title',
                    vmsElementsCount: vmsElements.length,
                    templatesElementsCount: templatesElements.length,
                    url: window.location.href
                };
            }''')
            print(f"主机模版页面内容: {json.dumps(templates_page_content, indent=2, ensure_ascii=False)}")
            results["dom_state"].append({
                "phase": "router.push 后",
                "state": templates_page_content
            })

            # 检查 Vue Router 状态
            router_state_after = page.evaluate('''() => {
                const app = document.querySelector('#app')?.__vue_app__;
                if (app) {
                    const router = app.config.globalProperties.$router;
                    return {
                        currentRoute: router.currentRoute.value?.path,
                        matchedCount: router.currentRoute.value?.matched?.length,
                        matchedPaths: router.currentRoute.value?.matched?.map(m => m.path)
                    };
                }
                return { error: 'Vue app not found' };
            }''')
            print(f"Router 状态: {json.dumps(router_state_after, indent=2, ensure_ascii=False)}")
            results["vue_state"].append({
                "phase": "router.push 后",
                "state": router_state_after
            })

            page.screenshot(path='/tmp/routing_templates_push.png')

            # 4. 分析问题
            print("\n=== 问题分析 ===")
            if router_state_after.get("currentRoute") == "/compute/host-templates":
                print("Router 路径正确: /compute/host-templates")
                if templates_page_content.get("pageTitle") == "虚拟机管理":
                    print("问题确认: Router 路径正确但页面标题仍显示虚拟机管理")
                    print("原因: 子组件 router-view 未正确响应路由变化")
                elif templates_page_content.get("templatesContainerExists"):
                    print("页面内容正确渲染")
                else:
                    print("页面内容异常，需要进一步分析")
            else:
                print(f"Router 路径异常: {router_state_after.get('currentRoute')}")

            # 保存报告
            with open('/tmp/routing_detailed_report.json', 'w') as f:
                json.dump(results, f, indent=2, ensure_ascii=False)

            print(f"\n报告已保存到 /tmp/routing_detailed_report.json")

        except Exception as e:
            print(f"错误: {e}")
            results["error"] = str(e)
            page.screenshot(path='/tmp/routing_error.png')

        finally:
            browser.close()

        return results

if __name__ == "__main__":
    diagnose()