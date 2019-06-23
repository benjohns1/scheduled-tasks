import { createUUIDs, createUUID } from '../../support/uuid';

describe('new schedule functionality', () => {
	
	describe('new schedule button', () => {
		it('creates an editable schedule form at the top', () => {
			cy.visitWait('/schedule');
			cy.get('[data-test=schedules]').then($t => $t.find('[data-test=schedule-item]').length).then(startingCount => {
				cy.get('[data-test=new-schedule-button]').click();
				const expectedCount = startingCount + 1;
				cy.get('[data-test=schedule-item]').should('have.length', expectedCount);
				cy.get('[data-test=schedule-item]').first().then($s => {
					cy.log('form inputs exist have expected default values');
					cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour');
					cy.wrap($s).find('[data-test=schedule-frequency-input]').should('have.value', 'Hour');
					cy.wrap($s).find('[data-test=schedule-interval-input]').should('have.value', '1');
					cy.wrap($s).find('[data-test=schedule-offset-input]').should('have.value', '0');
					cy.wrap($s).find('[data-test=schedule-at-minutes-input]').should('have.value', '0');
					cy.wrap($s).contains('[data-test=save-button]', 'save').click();
	
					cy.log('save button should make form input uneditable');
					cy.get('[data-test=schedule-item]').should('have.length', expectedCount);
					cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour');
					cy.wrap($s).find('[data-test=schedule-frequency]').should('have.text', 'Hour');
					cy.wrap($s).find('[data-test=schedule-interval]').should('have.text', '1');
					cy.wrap($s).find('[data-test=schedule-offset]').should('have.text', '0');
					cy.wrap($s).find('[data-test=schedule-at-minutes]').should('have.text', '0');
	
					cy.log('data persists after page reload');
					cy.visitWait('/schedule');
					cy.get('[data-test=schedule-item]').should('have.length', expectedCount);
					cy.get('[data-test=schedule-item]').first().then($rs => {
						cy.wrap($rs).contains('[data-test=open-button]', '>').click();
						cy.wrap($rs).contains('[data-test=close-button]', 'v');
						cy.wrap($rs).find('[data-test=schedule-name]').should('have.text', 'every hour');
						cy.wrap($rs).find('[data-test=schedule-frequency]').should('have.text', 'Hour');
						cy.wrap($rs).find('[data-test=schedule-interval]').should('have.text', '1');
						cy.wrap($rs).find('[data-test=schedule-offset]').should('have.text', '0');
						cy.wrap($rs).find('[data-test=schedule-at-minutes]').should('have.text', '0');
					});
				});
			});
		});

		it('saves custom schedule data for multiple schedules', () => {
			cy.visitWait('/schedule');
			cy.get('[data-test=schedules]').then($t => $t.find('[data-test=schedule-item]').length).then(startingCount => {
				let count = startingCount;
				const scheduleData = [
					{
						args: {
							frequency: 'Hour',
							interval: 1,
							offset: 0,
							atMinutes: '0',
						},
						want: {
							name: 'every hour',
							frequency: 'Hour',
							interval: '1',
							offset: '0',
							atMinutes: '0'
						}
					},
					{
						args: {
							frequency: 'Hour',
							interval: 2,
							offset: 1,
							atMinutes: '0,15',
						},
						want: {
							name: 'every 2 hours at 00, 15 minutes',
							frequency: 'Hour',
							interval: '2',
							offset: '1',
							atMinutes: '0,15'
						}
					},
					{
						args: {
							frequency: 'Hour',
							interval: 25,
							offset: 25,
							atMinutes: '30,0,61,5',
						},
						want: {
							name: 'every 24 hours at 00, 01, 05, 30 minutes',
							frequency: 'Hour',
							interval: '24',
							offset: '24',
							atMinutes: '0,1,5,30'
						}
					}
				];

				cy.log('add multiple schedules without saving')
				scheduleData.forEach(s => {
					cy.addSchedule(s.args, { save: false, visit: false });
					count++;
					cy.get('[data-test=schedule-item]').should('have.length', count);
				});

				cy.log('save schedules from the top down')
				const bottomEditSchedule = count - startingCount;
				for (let i = 1; i <= bottomEditSchedule; i++) {
					cy.get(`[data-test=schedule-item]:nth-child(${i})`).then($s => {
						cy.wrap($s).contains('[data-test=save-button]', 'save').click();
					});
				}

				cy.log('save button should make form fields uneditable');
				scheduleData.forEach((s, i) => {
					cy.get(`[data-test=schedule-item]:nth-child(${i+1})`).then($s => {
						cy.wrap($s).find('[data-test=schedule-name]').should('have.text', s.want.name);
						cy.wrap($s).find('[data-test=schedule-frequency]').should('have.text', s.want.frequency);
						cy.wrap($s).find('[data-test=schedule-interval]').should('have.text', s.want.interval);
						cy.wrap($s).find('[data-test=schedule-offset]').should('have.text', s.want.offset);
						cy.wrap($s).find('[data-test=schedule-at-minutes]').should('have.text', s.want.atMinutes);
					});
				});
				
				cy.log('data persists after page reload');
				cy.visitWait('/schedule');
				cy.get('[data-test=schedule-item]').should('have.length', count);
				scheduleData.forEach((s, i) => {
					cy.get(`[data-test=schedule-item]:nth-child(${i+1})`).then($s => {
						cy.wrap($s).contains('[data-test=open-button]', '>').click();
						cy.wrap($s).find('[data-test=schedule-name]').should('have.text', s.want.name);
						cy.wrap($s).find('[data-test=schedule-frequency]').should('have.text', s.want.frequency);
						cy.wrap($s).find('[data-test=schedule-interval]').should('have.text', s.want.interval);
						cy.wrap($s).find('[data-test=schedule-offset]').should('have.text', s.want.offset);
						cy.wrap($s).find('[data-test=schedule-at-minutes]').should('have.text', s.want.atMinutes);
					});
				});
			});
		});
	});
	
	describe('add task button', () => {
		it(`adds recurring tasks to a new schedule`, () => {
			const tasks = createUUIDs(2).map((id, index) => {
				return {
					name: `recurring task ${index}: ${id}`,
					description: `recurring task ${index} description: ${id}`
				};
			});
			cy.addSchedule({
				frequency: 'Hour',
				interval: 1,
				offset: 0,
				atMinutes: '0,30'
			}, {save: false});

			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour at 00, 30 minutes');

				cy.log('add 2 recurring tasks');
				cy.wrap($s).find('[data-test=new-task]').click();
				cy.wrap($s).find('[data-test=task-item]:nth-child(1)').then($ti => {
					cy.wrap($ti).find('[data-test=task-name-input]').clear().type(tasks[0].name);
					cy.wrap($ti).find('[data-test=task-description-input]').clear().type(tasks[0].description);
					cy.wrap($ti).find('[data-test=save-button]').should('not.exist');
				});
				
				cy.wrap($s).find('[data-test=new-task]').click();
				cy.wrap($s).find('[data-test=task-item]').then($tis => {
					cy.wrap($tis[0]).find('[data-test=task-name-input]').clear().type(tasks[1].name);
					cy.wrap($tis[0]).find('[data-test=task-description-input]').clear().type(tasks[1].description);
					cy.wrap($tis[0]).find('[data-test=save-button]').should('not.exist');
				});

				cy.log('save entire schedule with tasks');
				cy.wrap($s).find('[data-test=save-button]').click();
				cy.wrap($s).find('[data-test=task-item]').then($tis => {
					cy.wrap($tis.length).should('eq', tasks.length);
					cy.wrap($tis).find('[data-test=save-button]').should('not.exist');
					cy.wrap($tis[0]).find('[data-test=open-button]').click();
					cy.wrap($tis[1]).find('[data-test=open-button]').click();
					cy.wrap($tis).find('[data-test=task-name]').should('contain', tasks[0].name);
					cy.wrap($tis).find('[data-test=task-description]').should('contain', tasks[0].description);
					cy.wrap($tis).find('[data-test=task-name]').should('contain', tasks[1].name);
					cy.wrap($tis).find('[data-test=task-description]').should('contain', tasks[1].description);
				});
			});
			
			cy.log('ensure data persists after page reload');
			cy.visitWait('/schedule');
			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click();
				cy.wrap($s).find('[data-test=task-item]').then($tis => {
					cy.wrap($tis).find('[data-test=save-button]').should('not.exist');
					cy.wrap($tis[0]).find('[data-test=open-button]').click();
					cy.wrap($tis[1]).find('[data-test=open-button]').click();
					cy.wrap($tis).find('[data-test=task-name]').should('contain', tasks[0].name);
					cy.wrap($tis).find('[data-test=task-description]').should('contain', tasks[0].description);
					cy.wrap($tis).find('[data-test=task-name]').should('contain', tasks[1].name);
					cy.wrap($tis).find('[data-test=task-description]').should('contain', tasks[1].description);
				});
			});
		});

		it(`only adds first recurring task that is a duplicate`, () => {
			const id = createUUID();
			const duplicateTask = {
				name: `recurring task duplicate: ${id}`,
				description: `recurring task duplicate description: ${id}`
			}
			
			cy.addSchedule({
				frequency: 'Hour',
				interval: 1,
				offset: 0,
				atMinutes: '0,20'
			}, {save: false});

			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour at 00, 20 minutes');

				cy.log('add 2 recurring task duplicates');
				for (let i = 0; i < 2; i++) {
					cy.wrap($s).find('[data-test=new-task]').click();
					cy.wrap($s).find('[data-test=task-item]:nth-child(1)').then($ti => {
						cy.wrap($ti).find('[data-test=task-name-input]').clear().type(duplicateTask.name);
						cy.wrap($ti).find('[data-test=task-description-input]').clear().type(duplicateTask.description);
						cy.wrap($ti).find('[data-test=save-button]').should('not.exist');
					});
				}

				cy.log('save schedule with tasks');
				cy.wrap($s).find('[data-test=save-button]').click();
				cy.wrap($s).find('[data-test=task-item]').should('have.length', 1);
				cy.wrap($s).find('[data-test=task-item]').then($ti => {
					cy.wrap($ti).find('[data-test=save-button]').should('not.exist');
					cy.wrap($ti).find('[data-test=open-button]').click();
					cy.wrap($ti).find('[data-test=task-name]').should('have.text', duplicateTask.name);
					cy.wrap($ti).find('[data-test=task-description]').should('have.text', duplicateTask.description);
				});
			});
			
			cy.log('ensure data persists after page reload');
			cy.visitWait('/schedule');
			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click();
				cy.wrap($s).find('[data-test=task-item]').should('have.length', 1);
				cy.wrap($s).find('[data-test=task-item]').then($ti => {
					cy.wrap($ti).find('[data-test=save-button]').should('not.exist');
					cy.wrap($ti).find('[data-test=open-button]').click();
					cy.wrap($ti).find('[data-test=task-name]').should('have.text', duplicateTask.name);
					cy.wrap($ti).find('[data-test=task-description]').should('have.text', duplicateTask.description);
				});
			});
		});
	});
});