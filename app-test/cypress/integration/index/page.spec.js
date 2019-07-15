describe('index page basic elements', () => {
	beforeEach(() => {
		cy.visit('/')
	})

	it('has the expected title and page content', () => {
		cy.title().should('eq', 'Scheduled Tasks')
		cy.get('h1').should('have.text', 'Scheduled Tasks')
	})
})
