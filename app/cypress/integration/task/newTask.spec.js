const visitWait = url => {
	cy.visit(url).wait(500); // because of Sapper's script chunking, we need to wait extra for all svelte script chunks to be loaded after a page load
};

describe('new task functionality', () => {

	beforeEach(() => {
		visitWait('/task');
	});
	
	it('new task button creates an editable task form at the top', () => {
		cy.get('[data-test=task-item]').then($lis => {
			cy.get('[data-test=new-task-button]').click();
			cy.get('[data-test=task-item]').should('have.length', $lis.length + 1);
			cy.get('[data-test=task-item]').first().then($s => {
				cy.log('form inputs exist have expected default values');
				cy.wrap($s).find('[data-test=task-name-input]').should('have.value', 'new task');
				cy.wrap($s).find('[data-test=task-description-input]').should('have.value', '');
				cy.wrap($s).contains('[data-test=save-button]', 'save').click();

				cy.log('save button should make form input uneditable');
				cy.wrap($s).find('[data-test=task-name]').should('have.text', 'new task');
				cy.wrap($s).find('[data-test=task-description]').should('have.text', '');

				cy.log('data persists after page reload');
				visitWait('/task');
				cy.get('[data-test=task-item]').should('have.length', $lis.length + 1);
				cy.get('[data-test=task-item]').first().then($rs => {
					cy.wrap($rs).contains('[data-test=open-button]', '>').click();
					cy.wrap($rs).find('[data-test=task-name]').should('have.text', 'new task');
					cy.wrap($rs).find('[data-test=task-description]').should('have.text', '');
				});
			});
		});
	});
	
	it('save task saves custom task data', () => {
		cy.get('[data-test=new-task-button]').click();
		cy.get('[data-test=task-item]').first().then($s => {
			cy.wrap($s).find('[data-test=task-name-input]').clear().type('my test task name');
			cy.wrap($s).find('[data-test=task-description-input].description').clear().type('my test task description');
			cy.wrap($s).find('[data-test=save-button]').click();

			cy.log('save button saves custom task data');
			cy.wrap($s).find('[data-test=task-name]').should('have.text', 'my test task name');
			cy.wrap($s).find('[data-test=task-description]').should('have.text', 'my test task description');

			cy.log('custom task data persists after page reload');
			visitWait('/task');
			cy.get('[data-test=task-item]').first().then($rs => {
				cy.wrap($rs).find('[data-test=open-button]').click();
				cy.wrap($rs).find('[data-test=task-name]').should('have.text', 'my test task name');
				cy.wrap($rs).find('[data-test=task-description]').should('have.text', 'my test task description');
			});
		});
	});
});