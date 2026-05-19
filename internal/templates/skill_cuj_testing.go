package templates

const SkillCUJTestTemplate = `---
name: cuj-test
description: Run CUJ tests with headless Chromium and generate evidence reports
user-invocable: true
disable-model-invocation: true
---

# /cuj-test — Run CUJ tests with headless Chromium

## Trigger

User invokes ` + "`/cuj-test`" + ` to execute Critical User Journey tests against the running application.

## Instructions

You are running CUJ-based end-to-end tests using headless Chromium via Puppeteer. Each test follows a journey defined in ` + "`docs/CUJ.md`" + ` and produces evidence (screenshots, console logs, timing).

### Phase 1: Load CUJ Inventory

1. Read ` + "`docs/CUJ.md`" + ` for the full CUJ inventory
2. If the file doesn't exist, tell the user to run ` + "`/cuj-list`" + ` first to create the inventory
3. Parse all CUJ entries with their IDs, priorities, steps, and expected outcomes

### Phase 2: Select CUJs to Test

Ask the user which CUJs to test:

- **All** — run every CUJ in the inventory
- **By priority** — e.g., "P0 only" or "P0 and P1"
- **By ID** — e.g., "CUJ-001, CUJ-003, CUJ-007"
- **New/untested** — CUJs with no "Last Tested" date

Default to running all P0 CUJs if the user doesn't specify.

### Phase 3: Execute Tests

For each selected CUJ, run the test using headless Chromium via Puppeteer:

#### Setup

` + "```" + `javascript
const puppeteer = require('puppeteer');

const browser = await puppeteer.launch({
  headless: 'new',
  args: ['--no-sandbox', '--disable-setuid-sandbox']
});
const page = await browser.newPage();
await page.setViewport({ width: 1280, height: 720 });

// Monitor console output
const consoleLogs = [];
page.on('console', msg => {
  consoleLogs.push({ type: msg.type(), text: msg.text() });
});
const pageErrors = [];
page.on('pageerror', err => {
  pageErrors.push({ message: err.message, stack: err.stack });
});
` + "```" + `

#### For Each CUJ Step

1. **Execute the action** — navigate, type, click, select, etc.
2. **Wait for result** — use ` + "`waitForSelector`" + `, ` + "`waitForNavigation`" + `, or ` + "`waitForFunction`" + ` as appropriate
3. **Capture screenshot** — save to ` + "`screenshots/cuj-NNN-step-N-description.png`" + `
4. **Verify expected outcome** — check element presence, text content, URL changes
5. **Record timing** — track how long each step takes

` + "```" + `javascript
// Example step execution
const startTime = Date.now();

// Navigate
await page.goto('http://localhost:3000/signup', { waitUntil: 'networkidle0' });
await page.screenshot({ path: 'screenshots/cuj-001-step-1-signup-page.png', fullPage: true });

// Fill form
await page.type('#name', 'Test User');
await page.type('#email', 'test@example.com');
await page.type('#password', 'securepass123');
await page.screenshot({ path: 'screenshots/cuj-001-step-2-form-filled.png' });

// Submit and verify
await Promise.all([
  page.waitForNavigation({ waitUntil: 'networkidle0' }),
  page.click('button[type="submit"]')
]);

const successMsg = await page.$('.success-message');
if (!successMsg) {
  await page.screenshot({ path: 'screenshots/cuj-001-FAIL-step-3-no-success.png', fullPage: true });
  throw new Error('Expected success message not found');
}

const elapsed = Date.now() - startTime;
console.log('Step completed in ' + elapsed + 'ms');
` + "```" + `

#### Error Handling

- If a step fails, capture a failure screenshot immediately
- Record the exact error message and stack trace
- Continue to the next CUJ (don't abort the entire run)
- Mark the CUJ as FAIL with the failure details

### Phase 4: Generate Report

After all tests complete, produce a structured test report:

` + "```" + `markdown
# CUJ Test Report — YYYY-MM-DD

## Summary
- **Total CUJs tested:** N
- **Passed:** N
- **Failed:** N
- **Skipped:** N
- **Total duration:** Ns

## Results

### CUJ-001: User Registration — PASS (8.2s)
- Steps completed: 7/7
- Screenshots: cuj-001-step-1.png through cuj-001-step-7.png
- Console errors: 0

### CUJ-003: Checkout Flow — FAIL (12.1s)
- Failed at Step 4: "Apply coupon code"
- Expected: Discount applied, total updated
- Actual: Coupon field not responding to input
- Screenshot: cuj-003-FAIL-step-4.png
- Console errors: 1
  - TypeError: Cannot read property 'discount' of undefined
- **Reproduction steps:**
  1. Navigate to /cart
  2. Add any item
  3. Click "Apply Coupon"
  4. Type "SAVE20" — field does not accept input
` + "```" + `

### Phase 5: Update CUJ File

1. Update each tested CUJ's "Last Tested" date and "Result" in ` + "`docs/CUJ.md`" + `
2. Update the summary table at the top of the file
3. Stage and commit ` + "`docs/CUJ.md`" + ` with:
   ` + "```" + `
   Update CUJ test results: N passed, N failed
   ` + "```" + `

### Phase 6: Cleanup

After testing and reporting are complete, clean up all test artifacts:

#### Browser Cleanup
- Close all browser pages and the browser instance:
  ` + "```" + `javascript
  await page.close();
  await browser.close();
  ` + "```" + `
- Verify no orphaned Chromium processes remain

#### Artifact Management
- Ask the user whether to keep or remove screenshots from ` + "`screenshots/`" + `
  - **Keep** (default) — useful for review and evidence
  - **Remove** — clean slate for the next test run: ` + "`rm -rf screenshots/cuj-*`" + `
- Remove any temporary files created during testing

#### Test Data Cleanup
- If tests created any test data in the application (test users, test records, test orders, etc.), reverse those changes:
  - Delete test user accounts
  - Remove test records from databases
  - Undo any state mutations made during testing
- Document what was cleaned up and what was left in place

#### State Reset
- If the application state was modified during testing, reset it:
  - Clear test sessions/cookies
  - Reset any feature flags toggled during testing
  - Restore original configuration if modified
- Confirm the application is in a clean state for the next test run

#### Cleanup Report
- Summarize what was cleaned up:
  - Browser instances closed
  - Screenshots kept/removed (count and size)
  - Test data removed
  - State changes reversed

## Tips

- Ensure the application is running before starting tests — check the base URL first
- Create the ` + "`screenshots/`" + ` directory if it doesn't exist
- Use ` + "`waitUntil: 'networkidle0'`" + ` for pages that load async data
- Set reasonable timeouts (5–10 seconds) — if a page takes longer, that's a performance issue
- For multi-page flows, track cookies and session state across navigations
- If Puppeteer isn't installed, guide the user through ` + "`npx puppeteer`" + ` or ` + "`npm install puppeteer`" + `

## Project Context

- **Project:** {{.ProjectName}}
- **Owner:** {{.OwnerName}} ({{.OwnerEmail}})
- **CUJ File:** ` + "`docs/CUJ.md`" + `
- **Screenshots:** ` + "`screenshots/`" + ` directory
`
