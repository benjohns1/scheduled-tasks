const baseUrl = Cypress.config().baseUrl;
const startPaths = ['/', '/task', '/schedule'];
const navItems = [
  {
    testid: 'nav-home',
    label: 'home',
    endPath: '/'
  },
  {
    testid: 'nav-task',
    label: 'tasks',
    endPath: '/task'
  },
  {
    testid: 'nav-schedule',
    label: 'schedules',
    endPath: '/schedule'
  },
];

describe('navigation menu from all routes', () => {
  startPaths.forEach(startPath => {
    navItems.forEach(n => testNavLink(startPath, n.testid, n.label, n.endPath));
  });
});

function testNavLink(startPath, testid, label, endPath) {
  it('navigates to ' + endPath, () => {
    cy.visit(startPath);
    cy.get(`[data-test=${testid}]`).should('have.text', label).click();
    cy.url().should('eq', baseUrl + endPath);
  });
}