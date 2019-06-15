describe('schedule page basic elements', () => {
	beforeEach(() => {
		cy.visit('/schedule');
	});

	it('has the expected title and page content', () => {
		cy.title().should('eq', 'Scheduled Tasks - Schedules');
		cy.get('h1').should('have.text', 'Schedules');
	});
});