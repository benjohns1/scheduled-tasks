describe('New task functionality', () => {

	beforeEach(() => {
		cy.visit('/task').wait(50); // for some reason, without this delay the 'new task' button doesn't trigger DOM updates consistently
	});

	it('new task button creates a new task list element', () => {
		cy.get('section.tasks ul li').then($lis => {
			cy.get('button').contains('new task').click();
			cy.get('section.tasks ul li').should('have.length', $lis.length + 1);
		});
	});
	
	it('new task button has editable task name input field', () => {
		cy.get('button').contains('new task').click();
		cy.get('section.tasks section').first().get('header h2 input').should('have.value', 'new task');
	});

	it('new task button has editable task description textarea', () => {
		cy.get('button').contains('new task').click();
		cy.get('section.tasks section').first().get('.panel textarea').should('have.value', '');
	});

	it('new task form has save button', () => {
		cy.get('button').contains('new task').click();
		cy.get('section.tasks section').first().contains('button', 'save');
	});

	it('save task makes task name uneditable', () => {
		cy.get('button').contains('new task').click();
		cy.get('section.tasks section').first().contains('button', 'save').click();
		cy.get('section.tasks section').first().get('header h2').should('contain', 'new task');
	});

	it('save task makes task description uneditable', () => {
		cy.get('button').contains('new task').click();
		cy.get('section.tasks section').first().contains('button', 'save').click();
		cy.get('section.tasks section').first().get('.panel .description').should('contain', '');
	});

	it('save task persists new task after page reload', () => {
		cy.get('section.tasks ul li').then($lis => {
			cy.get('button').contains('new task').click();
			cy.get('section.tasks section').first().contains('button', 'save').click();
			cy.visit('/task').wait(10);
			cy.get('section.tasks ul li').should('have.length', $lis.length + 1);
			cy.get('section.tasks section').first().get('header h2').should('contain', 'new task');
		});
	});
});