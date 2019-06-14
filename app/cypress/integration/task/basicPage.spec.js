describe('task page basic elements', () => {
	beforeEach(() => {
		cy.visit('/task');
	});

	it('has the expected title and page content', () => {
		cy.title().should('eq', 'Scheduled Tasks - Tasks');
		cy.contains('h1', 'Tasks');
		cy.contains('[data-test=new-task-button]', 'new task');
		cy.contains('h1', 'Completed');
	});
});