import { createUUIDs } from '../../support/uuid';

describe('edit schedule functionality', () => {
		
	describe('add task button', () => {
		it(`adds recurring tasks to an existing schedule`, () => {
			
			cy.addSchedule({
				frequency: 'Hour',
				interval: 1,
				offset: 0,
				atMinutes: '0,30'
			});

			const tasks = createUUIDs(3).map((id, index) => {
				return {
					name: `recurring task ${index}: ${id}`,
					description: `recurring task description ${index}: ${id}`,
				};
			});

			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour at 00, 30 minutes');

				cy.log('add and save each recurring task');
				tasks.forEach(task => {
					cy.wrap($s).find('[data-test=new-task]').click();
					cy.wrap($s).find('[data-test=task-item]:nth-child(1)').then($ti => {
						cy.wrap($ti).find('[data-test=task-name-input]').clear().type(task.name);
						cy.wrap($ti).find('[data-test=task-description-input]').clear().type(task.description);
						cy.wrap($ti).find('[data-test=save-button]').click();

						cy.wrap($ti).find('[data-test=task-name]').should('have.text', task.name);
						cy.wrap($ti).find('[data-test=task-description]').should('have.text', task.description);
					});
				});
			});
			
			cy.log('ensure data persists after page reload');
			cy.visitWait('/schedule');
			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click();
				cy.wrap($s).find('[data-test=task-item]').should('have.length', tasks.length).then($tis => {
					$tis.toArray().forEach($ti => {
						cy.wrap($ti).find('[data-test=open-button]').click();
					});
					tasks.forEach(task => {
						cy.wrap($tis).find('[data-test=task-name]').should('contain', task.name);
						cy.wrap($tis).find('[data-test=task-description]').should('contain', task.description);
					});
				});
			});
				
		});
	});
});