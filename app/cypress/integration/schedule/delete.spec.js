import { createUUIDs } from '../../support/uuid'

describe('delete schedule functionality', () => {

	beforeEach(() => {
		cy.visitWait('/schedule')
	})
	
	describe('delete schedule button', () => {

		it(`deletes a new unsaved schedule`, () => {
			cy.get('[data-test=schedules]').then($t => $t.find('[data-test=schedule-item]').length).then(startingCount => {

				const tasks = createUUIDs(2).map((id, index) => {
					return {
						name: `delete schedule task ${index}: ${id}`,
						description: `delete schedule task description ${index}: ${id}`,
					}
				})
				
				cy.addSchedule({
					frequency: 'Hour',
					interval: 1,
					offset: 0,
					atMinutes: '0,28',
					tasks: tasks
				}, {save: false, visit: false})

				cy.get('[data-test=schedule-item]').should('have.length', startingCount + 1)
				cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
					cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour at 00, 28 minutes')
					cy.wrap($s).find('[data-test=task-item]:nth-child(1) [data-test=task-name-input]').should('have.value', tasks[1].name)
					cy.wrap($s).find('[data-test=task-item]:nth-child(2) [data-test=task-name-input]').should('have.value', tasks[0].name)
					cy.wrap($s).find('[data-test=delete-schedule-button]').click()
				})
				
				cy.get('[data-test=schedule-item]').should('have.length', startingCount)
				
			})
		})

		it(`deletes an existing schedule`, () => {
			cy.get('[data-test=schedules]').then($t => $t.find('[data-test=schedule-item]').length).then(startingCount => {

				const tasks = createUUIDs(2).map((id, index) => {
					return {
						name: `delete schedule task ${index}: ${id}`,
						description: `delete schedule task description ${index}: ${id}`,
					}
				})
				
				cy.addSchedule({
					frequency: 'Hour',
					interval: 1,
					offset: 0,
					atMinutes: '0,27',
					tasks: tasks
				}, {visit: false})

				cy.get('[data-test=schedule-item]').should('have.length', startingCount + 1)
				cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
					cy.wrap($s).find('[data-test=schedule-name]').should('have.text', 'every hour at 00, 27 minutes')
					cy.wrap($s).find('[data-test=task-item] [data-test=task-name]').should('contain', tasks[0].name)
					cy.wrap($s).find('[data-test=task-item] [data-test=task-name]').should('contain', tasks[1].name)
					cy.wrap($s).find('[data-test=delete-schedule-button]').click()
				})
				
				cy.get('[data-test=schedule-item]').should('have.length', startingCount)
				
				cy.log('ensure data persists after page reload')
				cy.visitWait('/schedule')
				cy.get('[data-test=schedule-item]').should('have.length', startingCount)

			})
		})

	})
})
