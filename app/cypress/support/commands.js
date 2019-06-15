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
Cypress.Commands.add("visitWait", url => {
	cy.visit(url).wait(1000);
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