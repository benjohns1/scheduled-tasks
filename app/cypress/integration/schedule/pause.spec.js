describe('edit schedule functionality', () => {

	beforeEach(() => {
		cy.visitWait('/schedule')
	})
	
	describe('pause checkbox', () => {
		it(`starts a new schedule paused/unpaused and pauses/unpauses an existing schedule`, () => {
			
			cy.addSchedule({
				frequency: 'Hour',
				interval: 1,
				offset: 0,
				atMinutes: '0,30',
				paused: false
			}, {visit: false})
			cy.get('[data-test=schedule-item]:nth-child(1) [data-test=paused-toggle]').should('exist').should('not.be.checked')

			cy.addSchedule({
				frequency: 'Day',
				interval: 1,
				offset: 0,
				atMinutes: '0,30',
				atHours: 6,
				paused: true
			}, {visit: false})
			cy.get('[data-test=schedule-item]:nth-child(1) [data-test=paused-toggle]').should('be.checked')

			cy.log('ensure paused state persists after reload')
			cy.visitWait('/schedule')
			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click()
				cy.wrap($s).find('[data-test=paused-toggle]').should('be.checked')
			})
			cy.get('[data-test=schedule-item]:nth-child(2)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click()
				cy.wrap($s).find('[data-test=paused-toggle]').should('exist').should('not.be.checked')
			})

			cy.log('pause/unpause schedules and check persistence')
			cy.get('[data-test=schedule-item]:nth-child(1) [data-test=paused-toggle]').uncheck({force: true})
			cy.get('[data-test=schedule-item]:nth-child(2) [data-test=paused-toggle]').check({force: true})
			cy.visitWait('/schedule')
			cy.get('[data-test=schedule-item]:nth-child(1)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click()
				cy.wrap($s).find('[data-test=paused-toggle]').should('exist').should('not.be.checked')
			})
			cy.get('[data-test=schedule-item]:nth-child(2)').then($s => {
				cy.wrap($s).find('[data-test=open-button]').click()
				cy.wrap($s).find('[data-test=paused-toggle]').should('be.checked')
			})
		})
	})
})
