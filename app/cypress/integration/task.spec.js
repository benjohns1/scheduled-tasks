describe('Task page basic elements', () => {
	beforeEach(() => {
		cy.visit('/task')
	});

	it('has the correct title', () => {
		cy.title().should('eq', 'Scheduled Tasks - Tasks');
	})

	it('has the correct headings', () => {
		cy.contains('h1', 'Tasks');
		cy.contains('h1', 'Completed');
	});

	it('navigates to home', () => {
		cy.get('nav a').contains('home').click();
		cy.url().should('include', '/');
	});

	it('navigates to /schedule', () => {
		cy.get('nav a').contains('schedules').click();
		cy.url().should('include', '/schedule');
	});

	it('has a new task button', () => {
		cy.get('button').contains('new task');
	});

	it('new task button creates a new task section at the top', () => {
		cy.get('button').contains('new task').click();
		cy.get('section.tasks section').first().get('header h2').should('eq', 'new task');
	})
});