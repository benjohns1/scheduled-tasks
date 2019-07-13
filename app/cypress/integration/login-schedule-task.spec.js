describe('login schedule task primary flow', () => {
  it('should login dev user, create a schedule, and create a task', () => {
    cy.devLogin('/schedule')
    cy.get('[data-test=new-schedule-button]').click()
    cy.get('[data-test=schedules] li[data-test=schedule-item]:nth-child(1)').then($sli => {
      cy.wrap($sli).find('[data-test=schedule-frequency-input]').select('Day')
      cy.wrap($sli).find('[data-test=schedule-at-hours-input]').clear().type('6,12')
      cy.wrap($sli).find('[data-test=schedule-at-minutes-input]').clear().type('0,30').blur()
      cy.wrap($sli).find('[data-test=new-task]').click()
      cy.wrap($sli).find('li[data-test=task-item]:nth-child(1)').then($tli => {
        cy.wrap($tli).find('[data-test=task-name-input]').clear().type('test recurring task')
        cy.wrap($tli).find('[data-test=task-description-input]').clear().type('test task description').blur()
      })
      cy.wrap($sli).find('[data-test=save-button]').click()
    })
    cy.visitWait('/task')
    cy.get('[data-test=new-task-button]').click()
    cy.get('[data-test=tasks] li[data-test=task-item]:nth-child(1)').then($tli => {
      cy.wrap($tli).find('[data-test=task-name-input]').clear().type('task name')
      cy.wrap($tli).find('[data-test=task-description-input]').clear().type('task description')
      cy.wrap($tli).find('[data-test=save-button]').click()
      cy.wrap($tli).find('[data-test=complete-toggle]').click()
    })
    cy.get('[data-test=completed-tasks] li[data-test=completed-task-item]:nth-child(1)').then($cli => {
      cy.wrap($cli).find('[data-test=task-name]').should('have.text', 'task name')
    })
    cy.get('[data-test=clear-tasks-button]').click()
    cy.get('[data-test=completed-tasks] li[data-test=completed-task-item]:nth-child(1)').should('not.exist')
  })
})