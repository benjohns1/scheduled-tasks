describe('new schedul functionality', () => {
  
	beforeEach(() => {
		cy.visitWait('/schedule');
	});
	
	describe('new schedule button', () => {
		it('creates an editable schedule form at the top', () => {
			cy.get('[data-test=schedules]').then($t => $t.find('[data-test=schedule-item]').length).then(startingCount => {
				cy.get('[data-test=new-schedule-button]').click();
				const expectedCount = startingCount + 1;
				cy.get('[data-test=schedule-item]').should('have.length', expectedCount);
				cy.get('[data-test=schedule-item]').first().then($s => {
					cy.log('form inputs exist have expected default values');
					cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'new schedule [unsaved]');
					cy.wrap($s).contains('[data-test=save-button]', 'save').click();
	
					cy.log('save button should make form input uneditable');
					cy.get('[data-test=schedule-item]').should('have.length', expectedCount);
					cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'hourly schedule');
	
					cy.log('data persists after page reload');
					cy.visitWait('/schedule');
					cy.get('[data-test=schedule-item]').should('have.length', expectedCount);
					cy.get('[data-test=schedule-item]').first().then($rs => {
						cy.wrap($rs).contains('[data-test=open-button]', '>').click();
						cy.wrap($rs).contains('[data-test=close-button]', 'v');
						cy.wrap($rs).find('[data-test=schedule-name]').should('have.text', 'hourly schedule');
					});
				});
			});
		});
	});
})