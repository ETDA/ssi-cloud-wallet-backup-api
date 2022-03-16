import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('vcs', function (table) {
        table.string('id', 255).primary()
        table.string('cid', 255).notNullable().unique()
        table.string('credential_type', 255).notNullable()
        table.string('schema_type', 255).notNullable()
        table.dateTime('issuance_date').notNullable()
        table.string('issuer', 255).notNullable()
        table.string('holder', 255).notNullable()
        table.text('jwt').notNullable()
        table.text('vc').notNullable()
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('vcs')
}
