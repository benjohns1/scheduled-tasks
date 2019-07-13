
describe('new schedule api functionality', () => {

  let fetchPolyfill

  before(() => {
    const polyfillUrl = 'https://unpkg.com/unfetch/dist/unfetch.umd.js'
    cy.request(polyfillUrl).then(response => {
      fetchPolyfill = response.body
    })
  })

  beforeEach(() => {
    cy.devLogin('/schedule', {
      onBeforeLoad(win) {
        delete win.fetch
        win.eval(fetchPolyfill)
        win.fetch = win.unfetch
      }
    })
  })
	
  it("doesn't send unnecessary schedule fields", () => {
    cy.server()
    cy.route('POST', '/schedule.json').as('postSchedule')

    const scheduleData = [
      {
        setup: {
          frequency: 'Month',
          interval: 2,
          offset: 1,
          atMinutes: '30',
          atHours: '8',
          onDaysOfMonth: '1, 15',
        },
        args: {
          frequency: 'Hour',
          interval: 1,
          offset: 0,
          atMinutes: '0',
        },
        wantEqual: {
          frequency: 'Hour',
          interval: 1,
          offset: 0,
          atMinutes: [0],
        },
        doNotWant: [
          'atHours',
          'onDaysOfMonth',
          'onDaysOfWeek',
          'tasks',
        ]
      },
      {
        setup: {
          frequency: 'Month',
          interval: 2,
          offset: 1,
          atMinutes: '30',
          atHours: '8',
          onDaysOfMonth: '1, 15',
        },
        args: {
          frequency: 'Day',
          interval: 2,
          offset: 1,
          atMinutes: '0,20',
          atHours: '6,12',
        },
        wantEqual: {
          frequency: 'Day',
          interval: 2,
          offset: 1,
          atMinutes: [0,20],
          atHours: [6,12],
        },
        doNotWant: [
          'onDaysOfMonth',
          'onDaysOfWeek',
          'tasks',
        ]
      },
      {
        setup: {
          frequency: 'Month',
          interval: 2,
          offset: 1,
          atMinutes: '30',
          atHours: '8',
          onDaysOfMonth: '1, 15',
        },
        args: {
          frequency: 'Week',
          interval: 1,
          offset: 0,
          atMinutes: '0,20',
          atHours: '6,12',
          onDaysOfWeek: ['Monday','Wednesday'],
        },
        wantEqual: {
          frequency: 'Week',
          interval: 1,
          offset: 0,
          atMinutes: [0,20],
          atHours: [6,12],
          onDaysOfWeek: ['Monday','Wednesday'],
        },
        doNotWant: [
          'onDaysOfMonth',
          'tasks',
        ]
      },
    ]

    scheduleData.forEach(s => {
      cy.addSchedule(s.setup, { save: false, visit: false })
      cy.addSchedule(s.args, { visit: false })
      cy.wait('@postSchedule').its('request').then(({body}) => {
        for (let key in s.wantEqual) {
          if (s.wantEqual.hasOwnProperty(key)) {
            cy.wrap(body[key]).debug('request prop', key).should('deep.equal', s.wantEqual[key])
          }
        }
        s.doNotWant.forEach(key => {
          cy.wrap(body[key]).debug('request prop', key).should('not.exist')
        })
      })
    })
  })
})
