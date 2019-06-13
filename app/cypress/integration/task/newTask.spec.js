describe('New task functionality', () => {

	beforeEach(() => {
		cy.visit('/task').wait(50); // for some reason, without this delay the 'new task' button doesn't trigger DOM updates consistently
	});

	it('new task button creates a new task list element', () => {
		cy.get('section.tasks ul li').then($lis => {
			cy.contains('button', 'new task').click();
			cy.get('section.tasks ul li').should('have.length', $lis.length + 1);
		});
	});
	
	it('new task button has editable task name input field', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().find('header h2 input').should('have.value', 'new task');
	});

	it('new task button has editable task description textarea', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().find('.panel textarea').should('have.value', '');
	});

	it('new task form has save button', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().contains('button', 'save');
	});

	it('save task makes task name uneditable', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).contains('button', 'save').click();
			cy.wrap($s).find('header h2').should('have.text', 'new task');
		});
	});

	it('save task makes task description uneditable', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).contains('button', 'save').click();
			cy.wrap($s).find('.panel .description').should('have.text', '');
		});
	});

	it('save task persists new task after page reload', () => {
		cy.get('section.tasks ul li').then($lis => {
			cy.contains('button', 'new task').click();
			cy.get('section.tasks section').first().contains('button', 'save').click();
			cy.visit('/task').wait(10);
			cy.get('section.tasks ul li').should('have.length', $lis.length + 1);
			cy.get('section.tasks section').first().find('header h2').should('have.text', 'new task');
		});
	});
	
	it('save task saves custom task data', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).find('header h2 input').clear().type('my test task name');
			cy.wrap($s).find('.panel textarea.description').clear().type('my test task description');
			cy.wrap($s).contains('button', 'save').click();
		});
		
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).find('header h2').should('have.text', 'my test task name');
			cy.wrap($s).find('.panel .description').should('have.text', 'my test task description');
		});
	});

	it('save task persists custom task data after page reload', () => {
		cy.contains('button', 'new task').click();
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).find('header h2 input').clear().type('my persistent test task name');
			cy.wrap($s).find('.panel textarea.description').clear().type('my persistent test task description');
			cy.wrap($s).contains('button', 'save').click();
		});
		cy.visit('/task').wait(10);
		
		cy.get('section.tasks section').first().then($s => {
			cy.wrap($s).contains('header button', '>').click();
			cy.wrap($s).find('header h2').should('have.text', 'my persistent test task name');
			cy.wrap($s).find('.panel .description').should('have.text', 'my persistent test task description');
		});
	});
});