import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('wallets', function (table) {
        table.string('id', 255).primary()
        table.string('did_address', 255).notNullable().unique()
        table.dateTime('created_at').notNullable()
        table.dateTime('deleted_at')
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('wallets')
}
