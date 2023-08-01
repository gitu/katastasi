import {expect, test} from '@playwright/test';
import type {Environment, StatusInfo} from "../src/types";
import {Critical, OK, Warning} from "../src/types";

test('visits the app root url', async ({page}) => {
    await page.route('*/**/api/envs', async route => {
        const json: Environment[] = [
            {id: "a", name: "Environment A", statusPages: {"x-full": "Full Overview Service X", "x-short": "Only Service X"}},
            {id: "b", name: "Enivronment B", statusPages: {"x-with-db": "X with DB", "x-full": "Full Overview Service X", "x-short": "Only Service X"}},
        ];
        await route.fulfill({json});
    });

    await page.goto('/');
    await expect(page.locator('div.title')).toHaveText('Status Pages');
    await page.screenshot({path: `screenshots/start-page.png`});
})


test('simple failing page', async ({page}) => {
    await page.route('*/**/api/env/prod/status/x', async route => {
        const json: StatusInfo = {
            lastUpdate: 0,
            overallStatus: Critical,
            services: [
                {
                    id: "x",
                    name: "Service X",
                    status: OK,
                    lastUpdate: 0,
                    serviceComponents: [
                        {status: OK, name: "Component A", info: "Everything is fine"},
                        {status: OK, name: "Component B", info: "Everything is fine"},
                        {status: OK, name: "Component C", info: "Everything is fine"},
                    ]
                },
                {
                    id: "y",
                    name: "Service Y",
                    status: Warning,
                    lastUpdate: 0,
                    serviceComponents: [
                        {status: OK, name: "Component A", info: "Everything is fine"},
                        {status: Warning, name: "Component B", info: "Something is wrong"},
                        {status: OK, name: "Component C", info: "Everything is fine"},
                    ]
                },
                {
                    id: "z",
                    name: "Service Z",
                    status: Critical,
                    lastUpdate: 0,
                    serviceComponents: [
                        {status: OK, name: "Component A", info: "Everything is fine"},
                        {status: OK, name: "Component B", info: "Everything is fine"},
                        {status: Critical, name: "Component C", info: "Something is wrong"},
                    ]
                }
            ],
            name: "Full Overview Service X"
        }
        await route.fulfill({json});
    });

    await page.goto('/env/prod/status/x');
    await expect(page.locator('#status-page-title')).toHaveText('Full Overview Service X');
    await page.screenshot({path: `screenshots/failing-status-page.png`, fullPage: true});
})


