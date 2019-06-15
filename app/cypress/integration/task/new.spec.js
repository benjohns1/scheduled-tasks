import { createUUID } from '../../support/uuid';

describe('new task functionality', () => {

	beforeEach(() => {
		cy.visitWait('/task');
	});
	
	describe('new task button', () => {
		it('creates an editable task form at the top', () => {
			cy.get('[data-test=tasks]').then($t => $t.find('[data-test=task-item]').length).then(startingCount => {
				cy.get('[data-test=new-task-button]').click();
				const expectedCount = startingCount + 1;
				cy.get('[data-test=task-item]').should('have.length', expectedCount);
				cy.get('[data-test=task-item]').first().then($s => {
					cy.log('form inputs exist have expected default values');
					cy.wrap($s).find('[data-test=task-name-input]').should('have.value', 'new task');
					cy.wrap($s).find('[data-test=task-description-input]').should('have.value', '');
					cy.wrap($s).contains('[data-test=save-button]', 'save').click();
	
					cy.log('save button should make form input uneditable');
					cy.get('[data-test=task-item]').should('have.length', expectedCount);
					cy.wrap($s).find('[data-test=task-name]').should('have.text', 'new task');
					cy.wrap($s).find('[data-test=task-description]').should('have.text', '');
	
					cy.log('data persists after page reload');
					cy.visitWait('/task');
					cy.get('[data-test=task-item]').should('have.length', expectedCount);
					cy.get('[data-test=task-item]').first().then($rs => {
						cy.wrap($rs).contains('[data-test=open-button]', '>').click();
						cy.wrap($rs).contains('[data-test=close-button]', 'v');
						cy.wrap($rs).find('[data-test=task-name]').should('have.text', 'new task');
						cy.wrap($rs).find('[data-test=task-description]').should('have.text', '');
					});
				});
			});
		});
	});
	
	describe('save task button', () => {
		it('saves custom task data', () => {
			cy.get('[data-test=new-task-button]').click();
			cy.get('[data-test=task-item]').first().then($s => {
				const id = createUUID();
				const name = 'test task name ' + id;
				const description = 'test task description' + id;

				cy.wrap($s).find('[data-test=task-name-input]').clear().type(name);
				cy.wrap($s).find('[data-test=task-description-input]').clear().type(description);
				cy.wrap($s).find('[data-test=save-button]').click();

				cy.log('save button saves custom task data');
				cy.wrap($s).find('[data-test=task-name]').should('have.text', name);
				cy.wrap($s).find('[data-test=task-description]').should('have.text', description);

				cy.log('custom task data persists after page reload');
				cy.visitWait('/task');
				cy.get('[data-test=task-item]').first().then($rs => {
					cy.wrap($rs).find('[data-test=open-button]').click();
					cy.wrap($rs).find('[data-test=task-name]').should('have.text', name);
					cy.wrap($rs).find('[data-test=task-description]').should('have.text', description);
				});
			});
		});
	});
});