// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

Cypress.Commands.add("debug", {prevSubject: true}, (subject, desc, value) => {
  Cypress.log({
    consoleProps: () => {
      return {
        'Description': desc,
        'Value': value,
      }
    }
	})
	return subject
})

// Because of Sapper's script chunking, we need to wait extra for all Svelte 
// script chunks to be loaded after a page load before using Svelte functionality
Cypress.Commands.add("visitWait", (url, options) => {
	cy.visit(url, options)
	cy.get('[data-test=loaded]')
	cy.wait(500)
})

Cypress.Commands.add("addTask", (name, description) => {
	cy.visitWait('/task')
	cy.get('[data-test=new-task-button]').click()
	cy.get('[data-test=task-item]').first().then($s => {
		cy.wrap($s).find('[data-test=task-name-input]').clear().type(name)
		cy.wrap($s).find('[data-test=task-description-input]').clear().type(description)
		cy.wrap($s).find('[data-test=save-button]').click()
	})
})

Cypress.Commands.add("addSchedule", ({ frequency, interval, offset, atMinutes, atHours, onDaysOfWeek, onDaysOfMonth, paused, tasks}, { save = true, visit = true } = {}) => {
	if (visit) {
		cy.visitWait('/schedule')
	}
	cy.get('[data-test=new-schedule-button]').click()
	cy.get('[data-test=schedule-item]').first().then($s => {
		cy.wrap($s).find('[data-test=schedule-frequency-input]').select(frequency)
		cy.wrap($s).find('[data-test=schedule-interval-input]').clear().type(interval)
		cy.wrap($s).find('[data-test=schedule-offset-input]').clear().type(offset)
		cy.wrap($s).find('[data-test=schedule-at-minutes-input]').clear().type(atMinutes).blur()
		if (frequency !== 'Hour') {
			cy.wrap($s).find('[data-test=schedule-at-hours-input]').clear().type(atHours).blur()
		}
		if (frequency === 'Month') {
			cy.wrap($s).find('[data-test=schedule-on-days-of-month-input]').clear().type(onDaysOfMonth).blur()
		}
		if (frequency === 'Week') {
			['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'].forEach(d => {
				if (onDaysOfWeek.includes(d)) {
					cy.wrap($s).find(`[data-test=schedule-on-days-of-week-input-${d}]`).check()
				} else {
					cy.wrap($s).find(`[data-test=schedule-on-days-of-week-input-${d}]`).uncheck()
				}
			})
		}
		if (paused) {
			cy.wrap($s).find('[data-test=paused-toggle]').check({force: true})
		}
		if (tasks) {
			cy.addRecurringTasks($s, tasks, {save: false})
		}
		if (save) {
			cy.wrap($s).find('[data-test=save-button]').click()
		}
	})
})

Cypress.Commands.add("addRecurringTasks", ($scheduleItem, tasks, { save = true } = {}) => {
	tasks.forEach(task => {
		cy.wrap($scheduleItem).find('[data-test=new-task]').click()
		cy.wrap($scheduleItem).find('[data-test=task-item]:nth-child(1)').then($ti => {
			cy.wrap($ti).find('[data-test=task-name-input]').clear().type(task.name)
			cy.wrap($ti).find('[data-test=task-description-input]').clear().type(task.description)
			if (save) {
				cy.wrap($ti).find('[data-test=save-button]').click()
			}
		})
	})
})

Cypress.Commands.add("devLogin", (redirect, redirectOpts) => {
	cy.visitWait('/?devlogin')
	if (redirect !== undefined) {
		cy.visitWait(redirect, redirectOpts)
	}
})