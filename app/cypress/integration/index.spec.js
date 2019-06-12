describe('Index page basic elements', () => {
	beforeEach(() => {
		cy.visit('/')
	});

	it('has the correct title', () => {
		cy.title().should('eq', 'Scheduled Tasks');
	})

	it('has the correct headings', () => {
		cy.contains('h1', 'Scheduled Tasks')
	});

	it('navigates to /task', () => {
		cy.get('nav a').contains('tasks').click();
		cy.url().should('include', '/task');
	});

	it('navigates to /schedule', () => {
		cy.get('nav a').contains('schedules').click();
		cy.url().should('include', '/schedule');
	});
});