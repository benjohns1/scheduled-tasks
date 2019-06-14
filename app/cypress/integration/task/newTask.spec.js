const visitWait = url => {
	cy.visit(url).wait(100); // for some reason, without this delay after page load, buttons do not trigger DOM changes
};

describe('new task functionality', () => {

	beforeEach(() => {
		visitWait('/task');
	});
	
	it('new task button creates an editable task form at the top', () => {
		cy.get('section.tasks ul li').then($lis => {
			cy.contains('button', 'new task').click();
			cy.get('section.tasks ul li').should('have.length', $lis.length + 1);
			cy.get('section.tasks section').first().then($s => {
				cy.log('form inputs exist have expected default values');
				cy.wrap($s).find('header h2 input').should('have.value', 'new task');
				cy.wrap($s).find('.panel textarea').should('have.value', '');
				cy.wrap($s).contains('button', 'save').click();

				cy.log('save button should make form input uneditable');
				cy.wrap($s).find('header h2').should('have.text', 'new task');
				cy.wrap($s).find('.panel .description').should('have.text', '');

				cy.log('data persists after page reload');
				visitWait('/task');
				cy.get('section.tasks ul li').should('have.length', $lis.length + 1);
				cy.get('section.tasks section').first().then($rs => {
					cy.wrap($rs).contains('header button', '>').click();
					cy.wrap($rs).find('header h2').should('have.text', 'new task');
					cy.wrap($rs).find('.panel .description').should('have.text', '');
				});
			});
		});
	});
	
	it('save task saves custom task data', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).find('header h2 input').clear().type('my test task name');
			cy.wrap($s).find('.panel textarea.description').clear().type('my test task description');
			cy.wrap($s).contains('button', 'save').click();

			cy.log('save button saves custom task data');
			cy.wrap($s).find('header h2').should('have.text', 'my test task name');
			cy.wrap($s).find('.panel .description').should('have.text', 'my test task description');

			cy.log('custom task data persists after page reload');
			visitWait('/task');
			cy.get('section.tasks section').first().then($rs => {
				cy.wrap($rs).contains('header button', '>').click();
				cy.wrap($rs).find('header h2').should('have.text', 'my test task name');
				cy.wrap($rs).find('.panel .description').should('have.text', 'my test task description');
			});
		});
	});
});