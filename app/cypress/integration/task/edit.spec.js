import { createUUID } from '../../support/uuid';

describe('edit task functionality', () => {

	describe('complete task button', () => {
		it('completes an existing task and moves it to the top of the completed list', () => {
			const id = createUUID();
			const name = 'complete test task name ' + id;
			const description = 'complete test task description ' + id;
			cy.addTask(name, description);
			cy.get('[data-test=task-item]').first().then($ti => {
				cy.wrap($ti).contains('[data-test=complete-toggle]', 'done');
				cy.wrap($ti).find('[data-test=task-name]').should('have.text', name);
				cy.wrap($ti).find('[data-test=task-description]').should('have.text', description);
				cy.wrap($ti).find('[data-test=complete-toggle]').click();
			});
			
			cy.log('task should be moved to completed list');
			cy.get('[data-test=completed-task-item]').first().then($cti => {
				cy.wrap($cti).find('[data-test=complete-toggle]').should('not.exist');
				cy.wrap($cti).find('[data-test=task-name]').should('have.text', name);
				cy.wrap($cti).find('[data-test=task-description]').should('have.text', description);
			});

			cy.log('reload page to test persistence');
			cy.visitWait('/task');
			cy.get('[data-test=completed-task-item]').first().then($cti => {
				cy.wrap($cti).find('[data-test=complete-toggle]').should('not.exist');
				cy.wrap($cti).find('[data-test=task-name]').should('have.text', name);
				cy.wrap($cti).find('[data-test=open-button]').click();
				cy.wrap($cti).find('[data-test=task-description]').should('have.text', description);
			});
		});
	});
});