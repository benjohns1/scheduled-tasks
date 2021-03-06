import { createUUID } from '../../support/uuid'

describe('edit task functionality', () => {

	beforeEach(() => {
		cy.devLogin('/task')
	})

	describe('clear task button', () => {
		it('clears all completed tasks', () => {
			cy.addTask('clear task 1', 'clear task 1 description')
			cy.addTask('clear task 2', 'clear task 2 description')
			cy.get('[data-test=task-item]:nth-child(1) [data-test=complete-toggle]').click()
			cy.get('[data-test=task-item]:nth-child(1) [data-test=complete-toggle]').click()

			cy.log('completed tasks should be cleared')
			cy.get('[data-test=clear-tasks-button]').click()
			cy.get('[data-test=completed-tasks]').then($ct => $ct ? $ct.find('[data-test=completed-task-item]').length : 0).then(count => {
				cy.wrap(count).should('eq', 0)
			})
			cy.contains('[data-test=completed-success-message]', 'Cleared all completed tasks')
			cy.contains('[data-test=completed-empty-message]', 'No completed tasks')

			cy.log('reload page to test persistence')
			cy.visitWait('/task')
			cy.get('[data-test=completed-tasks]').then($ct => $ct ? $ct.find('[data-test=completed-task-item]').length : 0).then(count => {
				cy.wrap(count).should('eq', 0)
			})
		})
	})

	describe('complete task button', () => {
		it('completes an existing task and moves it to the top of the completed list', () => {
			const id = createUUID()
			const name = 'complete test task name ' + id
			const description = 'complete test task description ' + id
			cy.addTask(name, description)
			cy.get('[data-test=task-item]').first().then($ti => {
				cy.wrap($ti).find('[data-test=task-name]').should('have.text', name)
				cy.wrap($ti).find('[data-test=task-description]').should('have.text', description)
				cy.wrap($ti).find('[data-test=complete-toggle]').click()
			})
			
			cy.log('task should be moved to completed list')
			cy.get('[data-test=completed-task-item]').first().then($cti => {
				cy.wrap($cti).find('[data-test=complete-toggle]').should('not.exist')
				cy.wrap($cti).find('[data-test=task-name]').should('have.text', name)
				cy.wrap($cti).find('[data-test=open-button]').click()
				cy.wrap($cti).find('[data-test=task-description]').should('have.text', description)
			})

			cy.log('reload page to test persistence')
			cy.visitWait('/task')
			cy.get('[data-test=completed-task-item]').first().then($cti => {
				cy.wrap($cti).find('[data-test=complete-toggle]').should('not.exist')
				cy.wrap($cti).find('[data-test=task-name]').should('have.text', name)
				cy.wrap($cti).find('[data-test=open-button]').click()
				cy.wrap($cti).find('[data-test=task-description]').should('have.text', description)
			})
		})
	})

})
