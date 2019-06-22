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

// Because of Sapper's script chunking, we need to wait extra for all Svelte 
// script chunks to be loaded after a page load before using Svelte functionality
Cypress.Commands.add("visitWait", (url, options) => {
	cy.visit(url, options);
	cy.get('[data-test=loaded]');
});

Cypress.Commands.add("addTask", (name, description) => {
	cy.visitWait('/task');
	cy.get('[data-test=new-task-button]').click();
	cy.get('[data-test=task-item]').first().then($s => {
		cy.wrap($s).find('[data-test=task-name-input]').clear().type(name);
		cy.wrap($s).find('[data-test=task-description-input]').clear().type(description);
		cy.wrap($s).find('[data-test=save-button]').click();
	});
});

Cypress.Commands.add("addSchedule", ({ frequency, interval, offset, atMinutes, paused}, { save = true, visit = true } = {}) => {
	if (visit) {
		cy.visitWait('/schedule');
	}
	cy.get('[data-test=new-schedule-button]').click();
	cy.get('[data-test=schedule-item]').first().then($s => {
		cy.wrap($s).find('[data-test=schedule-frequency-input]').select(frequency);
		cy.wrap($s).find('[data-test=schedule-interval-input]').clear().type(interval);
		cy.wrap($s).find('[data-test=schedule-offset-input]').clear().type(offset);
		cy.wrap($s).find('[data-test=schedule-at-minutes-input]').clear().type(atMinutes).blur();
		if (paused) {
			cy.wrap($s).find('[data-test=paused-toggle]').check();
		}
		if (save) {
			cy.wrap($s).find('[data-test=save-button]').click();
		}
	});
});