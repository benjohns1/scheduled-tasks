describe('Schedule page basic elements', () => {
	beforeEach(() => {
		cy.visit('/schedule')
	});

	it('has the correct title', () => {
		cy.title().should('eq', 'Scheduled Tasks - Schedules');
	})

	it('has the correct headings', () => {
		cy.contains('h1', 'Schedules')
	});

	it('navigates to home', () => {
		cy.get('nav a').contains('home').click();
		cy.url().should('include', '/');
	});

	it('navigates to /task', () => {
		cy.get('nav a').contains('tasks').click();
		cy.url().should('include', '/task');
	});
});